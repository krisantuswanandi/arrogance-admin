import { render, Box, Text, useApp, useInput, useStdout } from "ink";
import { useEffect, useState } from "react";
import { initializeApp, type App } from "firebase-admin/app";
import { getAuth, UserRecord } from "firebase-admin/auth";

const ENTER_ALTERNATE_SCREEN = "\x1b[?1049h";
const EXIT_ALTERNATE_SCREEN = "\x1b[?1049l";

const App = () => {
  const { exit } = useApp();
  const { stdout } = useStdout();

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

  // Listen for keypresses
  useInput((input, key) => {
    if (input === "q" || (key.ctrl && input === "c")) {
      exit();
    }
  });

  return (
    <Box
      justifyContent="center"
      alignItems="center"
      width={stdout.columns}
      height={stdout.rows}
    >
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

// Enter alternate screen
process.stdout.write(ENTER_ALTERNATE_SCREEN);
process.on("exit", () => {
  process.stdout.write(EXIT_ALTERNATE_SCREEN);
});

render(<App />);
