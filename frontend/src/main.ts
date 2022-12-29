import { createApp } from "vue";
import LandingApp from "./app/LandingApp.vue";
import AboutPage from "./app/AboutPage.vue";
import "./index.css";

const landing = createApp(LandingApp);
const about = createApp(AboutPage);

landing.mount("#app");
about.mount("#about");
