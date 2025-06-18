import { render, Box, Text, useApp, useInput, useStdout } from "ink";

const ENTER_ALTERNATE_SCREEN = "\x1b[?1049h";
const EXIT_ALTERNATE_SCREEN = "\x1b[?1049l";

const App = () => {
  const { exit } = useApp();
  const { stdout } = useStdout();

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
  );
};

// Enter alternate screen
process.stdout.write(ENTER_ALTERNATE_SCREEN);
process.on("exit", () => {
  process.stdout.write(EXIT_ALTERNATE_SCREEN);
});

render(<App />);
