import "./Home.css"

function Home() {
    return (
        <div className="hone-h">
            <div className="home">
                <h1 className="h1-h">Stock Trading App</h1>
                <p className="home-text">This is an application where you can view the price and graph of the last trades of the selected stock</p>
                <a className="ref-h" href="http://localhost:3000/reg">Let's get started</a>
            </div>
        </div>
      );
}

export default Home;