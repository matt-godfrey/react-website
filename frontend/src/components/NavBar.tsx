import { getCurrentUserQuery } from "@/api/getCurrentUserQuery";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createRootRoute,
  Link,
  Outlet,
  useNavigate,
} from "@tanstack/react-router";
import { Button } from "./ui/button";
interface NavBarProps {}

export default function NavBar({}: NavBarProps) {
  const { data: user } = useQuery(getCurrentUserQuery());
  console.log(user);
  const API_URL = import.meta.env.VITE_API_URL;
  const endpoint = `${API_URL}/logout`;
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const handleLogout = async () => {
    const response = await fetch(endpoint, {
      method: "POST",
      credentials: "include",
    });

    queryClient.setQueryData(["me"], null);

    // if (!response.ok) {
    //   throw new Error("Logout failed");
    // }

    queryClient.invalidateQueries({
      queryKey: ["me"],
    });

    navigate({ to: "/" });
  };
  const isLoggedIn = !!user;
  console.log(user);
  console.log(isLoggedIn);
  return (
    <div className="p-4 m-2 flex gap-6">
      <Link to="/" className="[&.active]:font-bold">
        Home
      </Link>
      <Link to="/about" className="[&.active]:font-bold">
        About
      </Link>
      {isLoggedIn ? (
        <Button onClick={handleLogout}>Logout</Button>
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
