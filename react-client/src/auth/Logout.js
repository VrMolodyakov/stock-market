import {useNavigate} from "react-router-dom";
const Logout = () =>{
    const navigate = useNavigate();
    console.log("inside logout");
    if(localStorage.getItem("access_token") != null)
        localStorage.removeItem("access_token");
    navigate("/auth");
};
export default Logout;