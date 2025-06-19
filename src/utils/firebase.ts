import { initializeApp } from "firebase-admin/app";
import { getAuth } from "firebase-admin/auth";

let app = initializeApp();

export function initFirebase() {
  if (!app) {
    app = initializeApp();
  }
  return app;
}

export const getFirebaseApp = initFirebase;

export function getAuthInstance() {
  return getAuth(app);
}

export function fetchUsers(limit: number) {
  return getAuthInstance().listUsers(limit);
}
