import { Outlet } from "react-router-dom"

const Layout = () => {
    return (
        <main className="mainApp">
            <Outlet />
        </main>
    )
}

export default Layout