import { Box, Text, useInput } from "ink";

export const User = ({ goBack }: { goBack: () => void }) => {
  useInput((_, key) => {
    if (key.escape || key.delete) {
      goBack();
    }
  });

  return (
    <Box paddingLeft={2}>
      <Text>
        <Text bold>Esc</Text> to go back
      </Text>
    </Box>
  );
};
