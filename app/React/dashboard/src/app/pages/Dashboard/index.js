import React, { useState, useEffect } from "react";
import Results from "../../components/Results";

const baseURL = "wss://localhost:8000/dashboard";

const Dashboard = () => {
  const [results, setResults] = useState([]);

  useEffect(() => {
    // in order to handle self-signed certificates we need to turn off the validation
    process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
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
