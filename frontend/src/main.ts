import { createApp } from "vue";
import LandingApp from "./app/LandingApp.vue";
import AboutPage from "./app/AboutPage.vue";
import IndexExplorer from "./app/IndexExplorer.vue";
import "./index.css";

const landing = createApp(LandingApp);
const about = createApp(AboutPage);
const maildir = createApp(IndexExplorer);

landing.mount("#landing");
about.mount("#about");
maildir.mount("#maildir");
