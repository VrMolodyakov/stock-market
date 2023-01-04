import Login from './auth/LoginForm';
import Register from './auth/RegisterForm';
import Prices from './stock/Prices';
import useAuth from './routing/useAuth'
import Code from './stock/StockCode'
import RequierAuth from './routing/RequireAuth';
import Logout from "./auth/Logout"
import { BrowserRouter, Routes, Route,Navigate,useNavigate } from "react-router-dom"
import { Navbar,Nav,Container} from 'react-bootstrap';
import axios from "axios";
import Layout from './Layout';
import Home from "./home/Home"
import "./App.css"

function App() {
  const navigate = useNavigate();
  const {auth,setAuth} = useAuth();

  const onLogout = () =>{
    Logout(auth, setAuth)
    const axiosInstance = axios.create({
      withCredentials: true
   })
   axiosInstance.get("http://localhost:8080/api/auth/logout")
    .then((response) => {
        console.log(response)
    })
    .catch((error) => {
        console.log(error.config);
    });
    navigate("/auth");
  };

  return (
    <>
    <Navbar collapseOnSelect expand="lg" bg="blue" variant="white">
    <Navbar.Brand className = "Home"  href="/">Home</Navbar.Brand>
    <Navbar.Brand className = "Charts"  href="/price">Charts</Navbar.Brand> 
     
      <Navbar.Toggle aria-controls="responsive-navbar-nav" />
      <Navbar.Collapse id="responsive-navbar-nav">
        <Nav className="me-auto">
        </Nav>
        <Nav className="links">
          <Nav.Link className="navBarLink" href="/auth">Sign In</Nav.Link>
          <Nav.Link className="navBarLink" href="/reg">Sign up</Nav.Link>
          <Nav.Link className="navBarLink" onClick={onLogout} >Log out</Nav.Link>
        </Nav>
      </Navbar.Collapse>
      
    </Navbar>

    <div className="App">
     
      <Routes>
          <Route path="/" element={<Layout />}>
            <Route path="" element={<Home />} />
            <Route path="auth" element={<Login />} />
            <Route path="reg" element={<Register />} />
            <Route element = {<RequierAuth/>}>
                <Route path="/:slug" element={<Code/>} />
                <Route path="price" element={<Prices />} />              
            </Route>
          </Route>
      </Routes>
    </div>
   </>
  );
}

export default App;
