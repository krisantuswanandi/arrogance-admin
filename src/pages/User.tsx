import type { UserRecord } from "firebase-admin/auth";
import { Box, Text, useInput } from "ink";

export const User = ({
  user,
  goBack,
}: {
  user: UserRecord;
  goBack: () => void;
}) => {
  useInput((_, key) => {
    if (key.escape || key.delete) {
      goBack();
    }
  });

  return (
    <Box paddingLeft={2} flexDirection="column">
      <Box>
        <Text bold>{user.uid}</Text>
      </Box>
      <Box marginTop={1}>
        <Text bold color="magenta">
          Exercises
        </Text>
      </Box>
      <Box marginTop={1}>
        <Text>
          <Text bold>Esc</Text> to go back
        </Text>
      </Box>
    </Box>
  );
};
