import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import NavBar from "../components/NavBar";

const RootLayout = () => (
  <>
    <NavBar />

    <hr />
    <Outlet />
    <TanStackRouterDevtools />
  </>
);

export const Route = createRootRoute({ component: RootLayout });
