import React from "react";
import MaterialTable from "material-table";

const Results = ({ results }) => {
  return (
    <div style={{ maxWidth: "100%" }}>
      <MaterialTable
        columns={[
          { title: "Start number", field: "sportsmen_id" },
          { title: "Name", field: "sportsmen_name" },
          { title: "Time", field: "time_finish" },
        ]}
        data={results}
        title="Results"
      />
    </div>
  );
};

export default Results;
