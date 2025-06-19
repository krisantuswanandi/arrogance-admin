import { withFullScreen } from "fullscreen-ink";
import { App } from "./src/App";
import { initFirebase } from "./src/utils/firebase";

initFirebase();

withFullScreen(<App />).start();
