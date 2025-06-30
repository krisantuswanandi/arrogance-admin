import { Box, Text, useInput } from "ink";
import { useEffect, useState } from "react";
import {
  fetchExercisesByUser,
  fetchHistoriesByUser,
  fetchProfilesByUser,
  deleteUserAndAllData,
} from "../utils/firebase";
import type { UserRecord } from "firebase-admin/auth";
import type { Exercise, History, Profile } from "../utils/types";

const List = ({
  title,
  list,
}: {
  title: string;
  list: { id: string; name: string }[];
}) => {
  return (
    <>
      <Box marginTop={1}>
        <Text bold color="magenta">
          {title}
        </Text>
      </Box>
      <Box>
        <Box flexDirection="column">
          {list.map((item, index) => (
            <Box key={item.id}>
              <Text>
                {index + 1}. [<Text bold>{item.id}</Text>] {item.name}
              </Text>
            </Box>
          ))}
        </Box>
      </Box>
    </>
  );
};

export const User = ({
  user,
  goBack,
}: {
  user: UserRecord;
  goBack: () => void;
}) => {
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [histories, setHistories] = useState<History[]>([]);
  const [loading, setLoading] = useState(true);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    setLoading(true);
    Promise.all([
      fetchProfilesByUser(user.uid),
      fetchExercisesByUser(user.uid),
      fetchHistoriesByUser(user.uid),
    ])
      .then(([userProfiles, userExercises, userHistories]) => {
        setProfiles(userProfiles);
        setExercises(userExercises);
        setHistories(userHistories);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [user.uid]);

  useInput((input, key) => {
    if (key.escape || key.delete) {
      goBack();
    }
    if (key.ctrl && input === "d") {
      handleDeleteUser();
    }
  });

  const handleDeleteUser = async () => {
    if (deleting) return;

    setDeleting(true);
    try {
      await deleteUserAndAllData(user.uid);
      goBack(); // Go back to users list after successful deletion
    } catch (error) {
      console.error("Failed to delete user:", error);
      setDeleting(false);
    }
  };

  let content;

  if (deleting) {
    content = (
      <Box marginTop={1}>
        <Text color="red">Deleting user and all related data...</Text>
      </Box>
    );
  } else if (loading) {
    content = (
      <Box marginTop={1}>
        <Text>Loading...</Text>
      </Box>
    );
  } else if (
    profiles.length === 0 &&
    exercises.length === 0 &&
    histories.length === 0
  ) {
    content = (
      <Box marginTop={1}>
        <Text color="gray">No data found</Text>
      </Box>
    );
  } else {
    const profilesContent = profiles.length ? (
      <List title="Profiles" list={profiles} />
    ) : null;
    const exercisesContent = exercises.length ? (
      <List title="Exercises" list={exercises} />
    ) : null;
    const historiesContent = histories.length ? (
      <List
        title="Histories"
        list={histories.map((h) => ({ id: h.id, name: h.workout.name }))}
      />
    ) : null;

    content = (
      <Box flexDirection="column">
        {profilesContent}
        {exercisesContent}
        {historiesContent}
      </Box>
    );
  }

  return (
    <Box paddingLeft={2} flexDirection="column">
      <Box>
        <Text bold>{user.uid}</Text>
      </Box>
      {content}
      <Box marginTop={1}>
        <Text>
          <Text bold>Esc</Text> to go back | <Text bold>Ctrl+D</Text> to delete
          user
        </Text>
      </Box>
    </Box>
  );
};
