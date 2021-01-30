import React from "react";

const Results = ({ results }) => {
  return (
    <>
      {results &&
        results.map((result) => (
          <p key={result.id}>
            <strong>
              {result.id},{result.checkpoint_id},{result.sportsmen_id},
              {result.event_state_id}
            </strong>
          </p>
        ))}
    </>
  );
};

export default Results;
