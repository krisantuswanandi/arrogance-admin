import { Box, Text, useApp, useInput } from "ink";
import { useEffect, useState } from "react";
import { initializeApp, type App } from "firebase-admin/app";
import { getAuth, UserRecord } from "firebase-admin/auth";
import { withFullScreen } from "fullscreen-ink";

const App = () => {
  const { exit } = useApp();
  const [users, setUsers] = useState<UserRecord[]>([]);

  // Initialize Firebase Admin SDK
  useEffect(() => {
    const app = initializeApp();

    getAuth(app)
      .listUsers(5)
      .then((result) => {
        setUsers(result.users);
      });
  }, []);

  // Exit on "q" (ctrl+C handled by fullscreen-ink)
  useInput((input) => {
    if (input === "q") {
      exit();
    }
  });

  return (
    <Box justifyContent="center" alignItems="center" width="100%" height="100%">
      <Box flexDirection="column">
        <Text>
          Press{" "}
          <Text bold color="red">
            "q"{" "}
          </Text>
          or{" "}
          <Text bold color="red">
            ctrl+C{" "}
          </Text>
          to exit.
        </Text>
        <Box marginTop={1} flexDirection="column">
          {users.map((user) => (
            <Text key={user.uid}>{user.uid}</Text>
          ))}
        </Box>
      </Box>
    </Box>
  );
};

withFullScreen(<App />).start();
