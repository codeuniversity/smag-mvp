import React, { useState, useEffect } from "react";
import Graph from "react-graph-vis";
import neo4j from "neo4j-driver";
import H1 from "./../components/H1";
import {
  UserIdRequest,
  UserSearchServicePromiseClient
} from "../protofiles/usersearch_grpc_web_pb";

export default function Network({ profile, foo }) {
  const [graph, setGraph] = useState(undefined);

  const options = {
    layout: {
      hierarchical: false
    },
    interaction: {
      dragNodes: false,
      dragView: false,
      zoomView: false
    },
    nodes: {
      font: "16px Helvetica white",
      size: 30,
      borderWidth: 1,
      color: {
        border: "black"
      },
      brokenImage:
        "https://www.autodeskfusionlifecycle.com/app/uploads/2017/06/default-user.png"
    },
    edges: {
      color: "#ffffff"
    },
    height: "1080px"
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
                label: user.userName
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
              size: 60
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
    gettNodesAndEdges(897148).then(result => setGraph(result));
  }, []);

  console.log(graph);
  if (graph) {
    return (
      <div>
        <H1>See who is following you</H1>
        <Graph graph={graph} options={options} events={events} />;
      </div>
    );
  }
  return <h1></h1>;
}
