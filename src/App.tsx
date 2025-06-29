import { Box, Spacer, Text, useApp, useInput } from "ink";
import { Users } from "./pages/Users";
import { Footer } from "./components/Footer";
import { useState } from "react";
import { User } from "./pages/User";
import type { UserRecord } from "firebase-admin/auth";

export const App = () => {
  const { exit } = useApp();
  const [page, setPage] = useState("users");
  const [selectedUser, setSelectedUser] = useState<UserRecord | null>(null);

  // Exit on "q" (ctrl+C handled by fullscreen-ink)
  useInput((input) => {
    if (input === "q") {
      exit();
    }
  });

  function selectUser(user: UserRecord) {
    setPage("user");
    setSelectedUser(user);
  }

  function goHome() {
    setPage("users");
    setSelectedUser(null);
  }

  let content;

  switch (page) {
    case "users":
      content = (
        <>
          <Box paddingLeft={2}>
            <Text bold>Users</Text>
          </Box>
          <Users onSelectUser={selectUser} />
        </>
      );
      break;
    case "user":
      content = <User user={selectedUser!} goBack={goHome} />;
      break;
    default:
      content = <Text color="red">404 Not Found</Text>;
  }

  return (
    <Box padding={1} paddingBottom={0}>
      <Box flexDirection="column" width="100%">
        <Box paddingLeft={2} marginBottom={1}>
          <Text color="magenta" bold>
            Arrogance Admin
          </Text>
        </Box>
        {content}
        <Spacer />
        <Footer />
      </Box>
    </Box>
  );
};
