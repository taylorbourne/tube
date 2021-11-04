import { ArrowBackIcon, ArrowForwardIcon } from '@chakra-ui/icons'
import { Box, Button } from '@chakra-ui/react'

type Props = {
  isFinalStep: boolean;
  isFirstStep: boolean;
  isFormValid: boolean;
  nextDisabled: boolean;
  nextStep: () => void;
  prevDisabled: boolean;
  prevStep: () => void;
};
const StepButtons = ({
  isFinalStep,
  isFirstStep,
  isFormValid,
  nextDisabled,
  nextStep,
  prevStep,
  prevDisabled,
}: Props) => {
  return (
    <Box
      bgColor="rgba(0, 0, 0, 0.25)"
      display="flex"
      justifyContent={isFirstStep ? "flex-end" : "space-between"}
      p="3"
    >
      {!isFirstStep && (
        <Button
          colorScheme="teal"
          variant="ghost"
          onClick={prevStep}
          disabled={prevDisabled}
          leftIcon={<ArrowBackIcon />}
          size="lg"
        >
          Back
        </Button>
      )}
      {isFinalStep && isFormValid && (
        <Button colorScheme="teal" size="lg" type="submit">
          Finish Setup
        </Button>
      )}
      {!isFinalStep && isFormValid && (
        <Button
          colorScheme="teal"
          variant="ghost"
          onClick={nextStep}
          disabled={nextDisabled}
          rightIcon={<ArrowForwardIcon />}
          size="lg"
        >
          Next
        </Button>
      )}
    </Box>
  );
};
export default StepButtons;
