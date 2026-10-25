import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
    layout("./navbar.layout.tsx", [
        index("routes/home.tsx"),
        route("/register", "routes/register.tsx")
    ])
] satisfies RouteConfig;
