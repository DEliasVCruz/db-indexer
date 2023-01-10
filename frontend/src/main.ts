import { createApp } from "vue";
import LandingPage from "./routes/LandingPage.vue";
import AboutPage from "./routes/AboutPage.vue";
import ExplorerApp from "./routes/ExplorerApp.vue";
import "./main.css";

const landing = createApp(LandingPage);
const about = createApp(AboutPage);
const maildir = createApp(ExplorerApp);

landing.mount("#landing");
about.mount("#about");
maildir.mount("#maildir");
