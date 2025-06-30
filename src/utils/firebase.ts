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

export function fetchUsers(limit: number, pageToken?: string) {
  return getAuthInstance().listUsers(limit, pageToken);
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

export async function fetchRoutinesByUser(uid: string) {
  const firestore = getFirestoreInstance();
  const routinesRef = firestore.collection("routines");
  const snapshot = await routinesRef.where("uid", "==", uid).get();

  return snapshot.docs.map((doc) => ({
    id: doc.id,
    ...doc.data(),
    createdAt: doc.data().createdAt.toDate(),
    updatedAt: doc.data().updatedAt.toDate(),
  }));
}

export async function deleteUserAndAllData(uid: string) {
  const auth = getAuthInstance();
  const firestore = getFirestoreInstance();

  try {
    // Delete from all collections where uid matches
    const collections = ["exercises", "histories", "routines"];

    // Delete from main collections
    for (const collectionName of collections) {
      const snapshot = await firestore
        .collection(collectionName)
        .where("uid", "==", uid)
        .get();
      const batch = firestore.batch();
      snapshot.docs.forEach((doc) => {
        batch.delete(doc.ref);
      });
      if (!snapshot.empty) {
        await batch.commit();
      }
    }

    // Delete profiles and their records subcollections
    const profilesSnapshot = await firestore
      .collection("profiles")
      .where("uid", "==", uid)
      .get();
    for (const profileDoc of profilesSnapshot.docs) {
      // Delete records subcollection for each profile
      const recordsSnapshot = await profileDoc.ref.collection("records").get();
      const recordsBatch = firestore.batch();
      recordsSnapshot.docs.forEach((recordDoc) => {
        recordsBatch.delete(recordDoc.ref);
      });
      if (!recordsSnapshot.empty) {
        await recordsBatch.commit();
      }

      // Delete the profile document itself
      await profileDoc.ref.delete();
    }

    // Finally, delete the user from Firebase Auth
    await auth.deleteUser(uid);

    return true;
  } catch (error) {
    console.error("Error deleting user and data:", error);
    throw error;
  }
}
