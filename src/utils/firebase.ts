import { initializeApp } from "firebase-admin/app";
import { getAuth } from "firebase-admin/auth";
import { getFirestore } from "firebase-admin/firestore";
import type { Exercise, History, Profile } from "./types";

let app = initializeApp();

export function initFirebase() {
  if (!app) {
    app = initializeApp();
  }
  return app;
}

export const getFirebaseApp = initFirebase;

export function getAuthInstance() {
  const app = getFirebaseApp();
  return getAuth(app);
}

export function getFirestoreInstance() {
  const app = getFirebaseApp();
  return getFirestore(app);
}

export function fetchUsers(limit: number) {
  return getAuthInstance().listUsers(limit);
}

export async function fetchProfilesByUser(uid: string) {
  const firestore = getFirestoreInstance();
  const profilesRef = firestore.collection("profiles");
  const snapshot = await profilesRef.where("uid", "==", uid).get();

  return snapshot.docs.map((doc) => ({
    id: doc.id,
    ...doc.data(),
    createdAt: doc.data().createdAt.toDate(),
    updatedAt: doc.data().updatedAt.toDate(),
  })) as Profile[];
}

export async function fetchExercisesByUser(uid: string) {
  const firestore = getFirestoreInstance();
  const exercisesRef = firestore.collection("exercises");
  const snapshot = await exercisesRef.where("uid", "==", uid).get();

  return snapshot.docs.map((doc) => ({
    id: doc.id,
    ...doc.data(),
    createdAt: doc.data().createdAt.toDate(),
    updatedAt: doc.data().updatedAt.toDate(),
  })) as Exercise[];
}

export async function fetchHistoriesByUser(uid: string) {
  const firestore = getFirestoreInstance();
  const historiesRef = firestore.collection("histories");
  const snapshot = await historiesRef.where("uid", "==", uid).get();

  return snapshot.docs.map((doc) => ({
    id: doc.id,
    ...doc.data(),
    workout: {
      ...doc.data().workout,
      date: doc.data().workout.date.toDate(),
    },
    createdAt: doc.data().createdAt.toDate(),
    updatedAt: doc.data().updatedAt.toDate(),
  })) as History[];
}
