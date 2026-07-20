import { createFileRoute } from "@tanstack/react-router";
import {  } from "@tanstack/react-query";

export const Route = createFileRoute("/logout")({
  component: RouteComponent,
});

function RouteComponent() {
  // const API_URL = import.meta.env.VITE_API_URL;
  // const endpoint = `${API_URL}/logout`;

  return <div>Hello "/logout"!</div>;
}
