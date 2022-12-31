import { createApp } from "vue";
import LandingApp from "./app/LandingApp.vue";
import AboutPage from "./app/AboutPage.vue";
import DirectoryApp from "./app/DirectoryApp.vue";
import "./index.css";

const landing = createApp(LandingApp);
const about = createApp(AboutPage);
const maildir = createApp(DirectoryApp);

landing.mount("#landing");
about.mount("#about");
maildir.mount("#maildir");
