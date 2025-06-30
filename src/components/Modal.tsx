import { Box, Text, useInput } from "ink";
import { useEffect, useState } from "react";

interface ConfirmationModalProps {
  isOpen: boolean;
  title?: string;
  confirmText?: string;
  cancelText?: string;
  onConfirm: () => void;
  onCancel: () => void;
  destructive?: boolean;
  children?: React.ReactNode;
}

const Button = ({
  text,
  isSelected,
  variant,
}: {
  text: string;
  isSelected: boolean;
  variant?: string;
}) => {
  let color = "gray";
  switch (variant) {
    case "destructive":
      color = "red";
      break;
    case "success":
      color = "green";
      break;
  }
  return (
    <Text color={color} bold={isSelected} inverse={isSelected}>
      {text}
    </Text>
  );
};

export const Modal = ({
  isOpen,
  title,
  confirmText = "Yes",
  cancelText = "No",
  onConfirm,
  onCancel,
  children,
  destructive = false,
}: ConfirmationModalProps) => {
  const [selectedOption, setSelectedOption] = useState<"confirm" | "cancel">(
    "cancel"
  );

  // Reset to cancel when modal opens
  useEffect(() => {
    if (isOpen) {
      setSelectedOption("cancel");
    }
  }, [isOpen]);

  useInput(
    (input, key) => {
      if (!isOpen) return;

      if (
        key.leftArrow ||
        key.rightArrow ||
        input === "h" ||
        input === "l" ||
        key.tab
      ) {
        setSelectedOption(selectedOption === "confirm" ? "cancel" : "confirm");
      }

      if (key.return) {
        if (selectedOption === "confirm") {
          onConfirm();
        } else {
          onCancel();
        }
      }

      if (key.escape) {
        onCancel();
      }
    },
    { isActive: isOpen }
  );

  if (!isOpen) return null;

  return (
    <Box
      position="absolute"
      width="100%"
      height="100%"
      alignItems="center"
      justifyContent="center"
    >
      <Box
        borderStyle="round"
        borderColor="white"
        paddingX={2}
        paddingY={1}
        minWidth={50}
        alignItems="center"
        justifyContent="center"
        flexDirection="column"
      >
        {title && (
          <Box marginBottom={1}>
            <Text bold color={destructive ? "red" : "blue"}>
              {title}
            </Text>
          </Box>
        )}

        <Box marginBottom={2}>
          <Text>{children}</Text>
        </Box>

        <Box gap={4}>
          <Button
            text={confirmText}
            isSelected={selectedOption === "confirm"}
            variant="destructive"
          />
          <Button text={cancelText} isSelected={selectedOption === "cancel"} />
        </Box>
      </Box>
    </Box>
  );
};
