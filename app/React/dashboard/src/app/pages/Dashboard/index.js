import React, { useState, useEffect, useCallback } from "react";
import Results from "../../components/Results";

const baseURL = "wss://localhost:8000/dashboard";

const Dashboard = () => {
  const [results, setResults] = useState([]);
  const [ws, setWs] = useState(null);

  const updateResults = useCallback(
    (msg) => {
      const newResult = JSON.parse(msg.data);

      // Check if a new message received has the same result ID.
      const isNew = !results.find((p) => p.id === newResult.id);

      // If ID is different, then append the existing results with a new row.
      // If ID is same, then grab time_finish and update existing row with the same result ID.
      if (isNew) {
        setResults([...results, newResult]);
      } else {
        const newResults = results.map((p) =>
          p.id === newResult.id
            ? { ...p, time_finish: newResult.time_finish }
            : p
        );
        setResults(newResults);
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
      <h1>Sports events dashboard ({results.length} elements)</h1>
      {results && <Results results={results} />}
      {console.log(results)}
    </div>
  );
};

export default Dashboard;
