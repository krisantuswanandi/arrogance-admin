export type Profile = {
  id: string;
  name: string;
  uid: string;
  createdAt: Date;
  updatedAt: Date;
};

export type Exercise = {
  id: string;
  name: string;
  uid: string;
  createdAt: Date;
  updatedAt: Date;
};

export type History = {
  id: string;
  uid: string;
  workout: Workout;
  createdAt: Date;
  updatedAt: Date;
};

export type Workout = {
  name: string;
  date: Date;
};

export type Routine = {
  id: string;
  name: string;
  uid: string;
  createdAt: Date;
  updatedAt: Date;
};
