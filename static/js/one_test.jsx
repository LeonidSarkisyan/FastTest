const TEST_ID = Number(
    window.location.pathname.split("/")[window.location.pathname.split("/").length - 1]
)

let response = await axios.get("/tests/" + TEST_ID)
console.log(response)

const title = response.data.title
const mainTitle = response.data.title

const App = () => {

    return (
        <div>
            <h1>{{ title }}</h1>
        </div>
    );
};

ReactDOM.render(<App />, document.getElementById("app"));
