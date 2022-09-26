import {useNavigate} from "react-router-dom";

const instance = axios.create({
    baseURL: "http://localhost:8080",
    headers: {
      "Content-Type": "application/json",
    },
  });

const removeCookie = async () =>{
    return instance.get("/api/auth/logout", userData);
  }

const Logout = () =>{
    const navigate = useNavigate();
    console.log("inside logout");
    if(localStorage.getItem("access_token") != null)
        localStorage.removeItem("access_token");
    (async() => {
        const response = await removeCookie();
        console.log(response)
        
    })();
    navigate("/auth");
};
export default Logout;