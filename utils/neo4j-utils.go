package utils

// Neo4jConfig holds all necessary informations to connect to a neo4j database
type Neo4jConfig struct {
	Host     string
	Username string
	Password string
}

//getNeo4jConfig returns a initialized Neo4jConfig object by reading the values from env variables
func GetNeo4jConfig() *Neo4jConfig {
	return &Neo4jConfig{
		Host:     GetStringFromEnvWithDefault("NEO4J_HOST", "localhost"),
		Username: GetStringFromEnvWithDefault("NEO4J_USERNAME", "neo4j"),
		Password: GetStringFromEnvWithDefault("NEO4J_PASSWORD", ""),
	}
}
