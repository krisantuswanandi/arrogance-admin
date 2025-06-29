import { Box, Text, useInput } from "ink";
import { useEffect, useState } from "react";
import { format } from "date-fns";
import { fetchUsers } from "../utils/firebase";
import type { UserRecord } from "firebase-admin/auth";

export const Users = ({
  onSelectUser,
}: {
  onSelectUser: (uid: UserRecord) => void;
}) => {
  const [users, setUsers] = useState<UserRecord[]>([]);
  const [cursor, setCursor] = useState<number>(0);

  useEffect(() => {
    fetchUsers(5).then((result) => {
      setUsers(result.users);
    });
  }, []);

  useInput((input, key) => {
    if (input === "j" || key.downArrow) {
      setCursor((prev) => Math.min(prev + 1, users.length - 1));
    } else if (input === "k" || key.upArrow) {
      setCursor((prev) => Math.max(prev - 1, 0));
    } else if (key.return) {
      if (users[cursor]) onSelectUser(users[cursor]);
    }
  });

  function formatDate(date: string | null | undefined) {
    if (!date) return "-";

    return format(new Date(date), "d MMM, HH:mm");
  }

  return (
    <Box marginTop={1} flexDirection="column">
      {users.length ? (
        users.map((user, index) => (
          <Box key={user.uid} flexDirection="column" marginBottom={1}>
            <Text color={index === cursor ? "magenta" : ""} bold>
              {index === cursor ? ">" : " "} {user.uid}
            </Text>
            <Box paddingLeft={2} flexDirection="column">
              <Text>Created at: {formatDate(user.metadata.creationTime)}</Text>
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
  );
};
