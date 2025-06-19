import { Box, Spacer, Text, useApp, useInput } from "ink";
import { Users } from "./pages/Users";
import { Footer } from "./components/Footer";

export const App = () => {
  const { exit } = useApp();

  // Exit on "q" (ctrl+C handled by fullscreen-ink)
  useInput((input) => {
    if (input === "q") {
      exit();
    }
  });

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
        <Users />
        <Spacer />
        <Footer />
      </Box>
    </Box>
  );
};
