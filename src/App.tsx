import { Box, Spacer, Text, useApp, useInput } from "ink";
import { Users } from "./pages/Users";
import { Footer } from "./components/Footer";
import { useState } from "react";
import { User } from "./pages/User";

export const App = () => {
  const { exit } = useApp();
  const [page, setPage] = useState("users");

  // Exit on "q" (ctrl+C handled by fullscreen-ink)
  useInput((input) => {
    if (input === "q") {
      exit();
    }
  });

  function selectPage(newPage: string) {
    setPage(newPage);
  }

  return (
    <Box padding={1} paddingBottom={0}>
      <Box flexDirection="column" width="100%">
        <Box paddingLeft={2} marginBottom={1}>
          <Text color="magenta" bold>
            Arrogance Admin
          </Text>
        </Box>
        {page === "users" ? (
          <>
            <Box paddingLeft={2}>
              <Text bold>Users</Text>
            </Box>
            <Users onSelect={selectPage} />
          </>
        ) : (
          <User goBack={() => setPage("users")} />
        )}
        <Spacer />
        <Footer />
      </Box>
    </Box>
  );
};
