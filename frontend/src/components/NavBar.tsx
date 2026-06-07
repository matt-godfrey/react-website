import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
interface NavBarProps {}

export default function NavBar({}: NavBarProps) {
  const isLoggedIn = false;
  return (
    <div className="p-4 m-2 flex gap-6">
      <Link to="/" className="[&.active]:font-bold">
        Home
      </Link>
      <Link to="/about" className="[&.active]:font-bold">
        About
      </Link>
      {isLoggedIn ? (
        <button>Log Out</button>
      ) : (
        <>
          <Link to="/login" className="[&.active]:font-bold">
            Login
          </Link>
          {/*<Link to="/register">Register</Link>*/}
        </>
      )}
    </div>
  );
}
