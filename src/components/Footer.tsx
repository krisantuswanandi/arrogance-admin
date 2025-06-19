import { Box, Text } from "ink";

export const Footer = () => {
  return (
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
  );
};
