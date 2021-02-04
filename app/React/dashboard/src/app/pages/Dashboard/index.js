import React, { useState, useEffect, useCallback } from "react";
import Results from "../../components/Results";

const baseURL = "wss://localhost:8000/dashboard";

const Dashboard = () => {
  const [results, setResults] = useState([]);
  const [ws, setWs] = useState(null);

  const updateResults = useCallback(
    (msg) => {
      // Parse the results JSON response.
      const data = JSON.parse(msg.data);
      // If data is not null then use it.
      if (data) {
        // Format time when data arrives (new data on just opened page) or new message.
        for (var i = 0; i < data.length; i++) {
          // If result has time_finish (UNIX time) then convert it into readable time.
          if (data[i].time_finish) {
            var fulldate = new Date(data[i].time_finish);
            data[i].time_finish = `${fulldate.getHours()}:${fulldate.getMinutes()}:${fulldate.getSeconds()}.${fulldate.getMilliseconds()}`;
          } else {
            data[i].time_finish = "";
          }
        }

        // If there are no results on page yet, use all data from the server.
        if (!results.length) {
          // If it's initial data then it's array, if that's a message then it's object.
          if (Array.isArray(data)) {
            setResults(data);
          } else {
            setResults([data, ...results]);
          }
        }

        // If there are messages in client state already, append these with new data from message.
        if (results.length) {
          // Check if a new message has the same result ID.
          const isNew = !results.find((p) => p.id === data.id);

          // If ID is different, then append the existing results with a new row.
          // If ID is same, then grab time_finish and update existing row with the same result ID.
          if (isNew) {
            setResults([data, ...results]);
          } else {
            const newResults = results.map((p) =>
              p.id === data.id ? { ...p, time_finish: data.time_finish } : p
            );
            setResults(newResults);
          }
        }
      }
    },
    [results, setResults]
  );

  useEffect(() => {
    const ws = new WebSocket(baseURL);

    ws.onopen = (evt) => {
      console.log("WebSocket opened!", { evt });
    };

    ws.onclose = (evt) => {
      console.log("WebSocket closed!", { evt });
      ws.current = undefined;
    };

    ws.onerror = (err) => {
      console.log("Websocket error:", { err });
    };

    setWs(ws);

    return () => {
      ws.close();
    };
  }, []);

  useEffect(() => {
    if (ws) {
      ws.onmessage = updateResults;
    }
  }, [ws, updateResults]);

  return (
    <div className="dashboard">
      <h1>Sport events dashboard</h1>
      {results && <Results results={results} />}
    </div>
  );
};

export default Dashboard;
