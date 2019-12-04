import React, { useState, useEffect } from "react";
import Graph from "react-graph-vis";
import neo4j from "neo4j-driver";
import { stringify } from "querystring";

export default function Network({ profile }) {
  const [graph, setGraph] = useState(undefined);

  const options = {
    layout: {
      hierarchical: true
    },
    edges: {
      color: "#ffffff"
    },
    height: "1000px"
  };

  const events = {
    select: function(event) {
      var { nodes, edges } = event;
    }
  };

  const gettNodesAndEdges = id => {
    const driver = neo4j.driver("bolt://localhost:7687", neo4j.auth.basic());
    const session = driver.session();

    const graphResult = {
      nodes: [],
      edges: []
    };

    graphResult.nodes.push({
      id: id,
      label: String(id),
      title: "USER node"
    });

    return session
      .run(`match p=()-[:FOLLOWS]->(:USER{id: ${id}}) return p`, {})
      .then(function(result) {
        console.log(result);
        result.records.forEach(element => {
          const start = element.get("p").start.properties.id.low;
          const end = element.get("p").end.properties.id.low;

          graphResult.nodes.push({
            id: start,
            label: String(start),
            title: "USER node"
          });
          graphResult.edges.push({ from: start, to: end });
        });
        session.close();
        driver.close();

        return graphResult;
      })
      .catch(function(err) {
        session.close();
        driver.close();
        console.log(err);
      });
  };

  useEffect(() => {
    gettNodesAndEdges(2282636).then(result => setGraph(result));
  }, []);

  console.log(graph);
  if (graph) {
    return <Graph graph={graph} options={options} events={events} />;
  }
  return <h1>Hallo</h1>;
}
