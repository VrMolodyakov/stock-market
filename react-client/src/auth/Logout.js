
const Logout = (auth, setAuth) => {
  const token = '';
  setAuth({token});
  localStorage.removeItem("access_token");
  localStorage.clear();
};

export default Logout;