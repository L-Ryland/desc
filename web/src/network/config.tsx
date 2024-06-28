const serverName = "test"

const config = {
    production: {
        apiUrl: `http://${ serverName }.dytx2tyxt.com/v1/`,
        debugMode: false,
    },
    development: {
        apiUrl: 'http://localhost:8071/v1/',
        debugMode: true,
    },
    test: {
        apiUrl: 'http://localhost:8071/v1/',
        debugMode: true,
    },
}

export default config
