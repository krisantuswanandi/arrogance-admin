import { Box, Spacer, Text, useApp, useInput } from "ink";
import { useEffect, useState } from "react";
import { initializeApp, type App } from "firebase-admin/app";
import { getAuth, UserRecord } from "firebase-admin/auth";
import { withFullScreen } from "fullscreen-ink";
import { format } from "date-fns";

const App = () => {
  const { exit } = useApp();
  const [users, setUsers] = useState<UserRecord[]>([]);
  const [cursor, setCursor] = useState<number>(0);

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
  useInput((input, key) => {
    if (input === "q") {
      exit();
    } else if (input === "j" || key.downArrow) {
      setCursor((prev) => Math.min(prev + 1, users.length - 1));
    } else if (input === "k" || key.upArrow) {
      setCursor((prev) => Math.max(prev - 1, 0));
    }
  });

  function formatDate(date: string | null | undefined) {
    if (!date) return "-";

    return format(new Date(date), "d MMM, HH:mm");
  }

  return (
    <Box padding={1} paddingBottom={0}>
      <Box flexDirection="column" width="100%">
        <Box paddingLeft={2}>
          <Text color="magenta" bold>
            Arrogance Admin
          </Text>
        </Box>
        <Box paddingLeft={2} marginTop={1}>
          <Text bold>Users</Text>
        </Box>
        <Box marginTop={1} flexDirection="column">
          {users.length ? (
            users.map((user, index) => (
              <Box key={user.uid} flexDirection="column" marginBottom={1}>
                <Text color={index === cursor ? "magenta" : ""} bold>
                  {index === cursor ? ">" : " "} {user.uid}
                </Text>
                <Box paddingLeft={2} flexDirection="column">
                  <Text>
                    Created at: {formatDate(user.metadata.creationTime)}
                  </Text>
                  <Text>
                    Last login: {formatDate(user.metadata.lastRefreshTime)}
                  </Text>
                </Box>
              </Box>
            ))
          ) : (
            <Box paddingLeft={2}>
              <Text italic>Loading...</Text>
            </Box>
          )}
        </Box>
        <Spacer />
        <Box paddingLeft={2}>
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
        </Box>
      </Box>
    </Box>
  );
};

withFullScreen(<App />).start();
