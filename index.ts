import { withFullScreen } from "fullscreen-ink";
import { createElement } from "react";
import { App } from "./src/App";
import { initFirebase } from "./src/utils/firebase";

initFirebase();

withFullScreen(createElement(App)).start();
