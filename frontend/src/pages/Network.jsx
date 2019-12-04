import React, { useState, useEffect } from "react";
import Graph from "react-graph-vis";
import neo4j from "neo4j-driver";
import {
  UserIdRequest,
  UserSearchServicePromiseClient
} from "../protofiles/usersearch_grpc_web_pb";

export default function Network({ profile, foo }) {
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
  const apiClient = new UserSearchServicePromiseClient("http://localhost:4000");

  const getUserData = async id => {
    const userIdRequest = new UserIdRequest();
    userIdRequest.setUserId(id);
    const response = await apiClient.getUserWithUserId(userIdRequest);
    return response;
  };

  const gettNodesAndEdges = id => {
    const driver = neo4j.driver("bolt://localhost:7687", neo4j.auth.basic());
    const session = driver.session();

    const graphResult = {
      nodes: [],
      edges: []
    };

    return session
      .run(`match p=()-[:FOLLOWS]->(:USER{id: ${id}}) return p`, {})
      .then(function(result) {
        console.log(result);
        return Promise.all(
          result.records.map(element => {
            const start = element.get("p").start.properties.id.low;
            const end = element.get("p").end.properties.id.low;

            return getUserData(start.toString()).then(protoUser => {
              const user = protoUser.toObject();
              graphResult.nodes.push({
                id: start,
                shape: "circularImage",
                image: user.avatarUrl,
                label: user.userName,
                title: "USER node"
              });
              graphResult.edges.push({ from: start, to: end });
            });
          }),
          getUserData(id.toString()).then(protoUser => {
            const user = protoUser.toObject();
            graphResult.nodes.push({
              id: id,
              shape: "circularImage",
              image: user.avatarUrl,
              label: user.userName,
              title: "USER node"
            });
          })
        ).then(() => {
          session.close();
          driver.close();

          return graphResult;
        });
      })
      .catch(function(err) {
        session.close();
        driver.close();
        console.log(err);
      });
  };

  useEffect(() => {
    gettNodesAndEdges(profile.id).then(result => setGraph(result));
  }, []);

  console.log(graph);
  if (graph) {
    return <Graph graph={graph} options={options} events={events} />;
  }
  return <h1></h1>;
}
