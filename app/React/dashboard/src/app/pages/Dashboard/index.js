import React, { useState, useEffect } from "react";
import Results from "../../components/Results";

const baseURL = "ws://localhost:8080/dashboard";

const Dashboard = () => {
  const [results, setResults] = useState([]);

  useEffect(() => {
    let ws = new WebSocket(baseURL);

    ws.onopen = (evt) => {
      console.log("WebSocket opened!", { evt });
    };

    ws.onclose = (evt) => {
      console.log("WebSocket closed!", { evt });
      ws = undefined;
    };

    ws.onmessage = (msg) => {
      console.log("Websocket message:", { msg });
      setResults(results.concat([JSON.parse(msg.data)]));
    };

    ws.onerror = (err) => {
      console.log("Websocket error:", { err });
    };
  }, []);

  return (
    <div className="dashboard">
      <h1>Sports events dashboard</h1>
      {results && <Results results={results} />}
    </div>
  );
};

export default Dashboard;
