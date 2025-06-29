import { Box, Text, useInput } from "ink";
import { useCallback, useEffect, useState } from "react";
import { format } from "date-fns";
import { fetchUsers } from "../utils/firebase";
import type { UserRecord } from "firebase-admin/auth";

export const Users = ({
  onSelectUser,
}: {
  onSelectUser: (uid: UserRecord) => void;
}) => {
  const [users, setUsers] = useState<UserRecord[]>([]);
  const [cursor, setCursor] = useState(0);
  const [currentPage, setCurrentPage] = useState(0);
  const [pageTokens, setPageTokens] = useState<string[]>([]);
  const [hasNextPage, setHasNextPage] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadPage();
  }, []);

  const loadPage = useCallback(
    async (pageToken?: string) => {
      setLoading(true);
      try {
        const result = await fetchUsers(3, pageToken);
        setUsers(result.users);
        setHasNextPage(!!result.pageToken);

        // Update page tokens for navigation
        if (pageToken && !pageTokens.includes(pageToken)) {
          setPageTokens((prev) => [...prev, pageToken]);
        }
      } finally {
        setLoading(false);
      }
    },
    [pageTokens]
  );

  useInput((input, key) => {
    if (loading) return; // Prevent input during loading

    if (input === "j" || key.downArrow) {
      setCursor((prev) => Math.min(prev + 1, users.length - 1));
    } else if (input === "k" || key.upArrow) {
      setCursor((prev) => Math.max(prev - 1, 0));
    } else if (key.return) {
      if (users[cursor]) onSelectUser(users[cursor]);
    } else if (input === "l" || key.rightArrow) {
      // Next page
      if (hasNextPage && users.length > 0) {
        const lastUser = users[users.length - 1];
        if (lastUser?.uid) {
          setCurrentPage((prev) => prev + 1);
          setCursor(0);
          loadPage(lastUser.uid);
        }
      }
    } else if (input === "h" || key.leftArrow) {
      // Previous page
      if (currentPage > 0) {
        const prevPage = currentPage - 1;
        const prevPageToken =
          prevPage === 0 ? undefined : pageTokens[prevPage - 1];
        setCurrentPage(prevPage);
        setCursor(0);
        loadPage(prevPageToken);
      }
    }
  });

  function formatDate(date: string | null | undefined) {
    if (!date) return "-";

    return format(new Date(date), "d MMM, HH:mm");
  }

  return (
    <Box marginTop={1} flexDirection="column">
      {loading ? (
        <Box paddingLeft={2} marginBottom={1}>
          <Text italic>Loading...</Text>
        </Box>
      ) : users.length ? (
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
          <Text italic>No users found</Text>
        </Box>
      )}
      <Box paddingLeft={2}>
        <Text>{currentPage > 0 && "← Prev | "}</Text>
        <Text bold>Page {currentPage + 1}</Text>
        <Text>{hasNextPage && " | Next →"}</Text>
      </Box>
    </Box>
  );
};
