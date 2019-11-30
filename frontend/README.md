## Social Record front-end

The Social Record front-end displays the analyzed data from instagram of the visitor of the exhibition. Using the front-end the visitors can experience what public data social media providers expose about them.

### Prerequisities

|                |                                                                    |     |     |     |
| -------------- | ------------------------------------------------------------------ | --- | --- | --- |
| Node.js        | `v10.16.0                                                          |
| Go             | `go 1.13` ([go modules](https://blog.golang.org/using-go-modules)) |     |     |     |
| Docker         | `v19.x`                                                            |     |     |     |
| Docker Compose | `v1.25.x`                                                          |     |     |     |

### Getting started

**1. Clone repo** [`https://github.com/codeuniversity/smag-mvp.git`](https://github.com/codeuniversity/smag-mvp.git)

**2. Switch to the front-end folder** `cd smag-mvp/frontend`

**3. Install all dependencies** `npm install`

**4. Run the application in development mode** `npm start` (runs on `localhost:3000`)

**5. To deploy to production you can create a minified bundle** `npm run build`

**6. Run all services in docker to locally test the application.**

1. Start all services `make run`
2. Add `127.0.0.1 my-kafka` and `127.0.0.1. minio` to your `/etc/hosts` file
3. Choose a user_name as a starting point and run `go run cli/main/main.go instagram <user_name>`

### React Component Design

| Name            | Prop Name    | Data Structure | Example               | Description                                                        |
| --------------- | ------------ | -------------- | --------------------- | ------------------------------------------------------------------ |
| `<Test/>`       | h1           | String         | `""`                  | Displays the h1 headline in the App component and Start component. |
| `<IGPost/>`     | key          | Function       | `(post.shortcode)`    | Displays the instagram post in the Result component.               |
| `<Form/>`       | onSubmit     | Function       | `(name.string)=>void` | Displays the search form in the App component.                     |
| `Camera-Feed/>` | onFileSubmit | Function       | `(file:File)=>void`   | Implements the camera feed in the Start component.                 |
