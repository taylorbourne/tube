import { Alert, AlertIcon, Box, Button, FormControl, FormErrorMessage, Heading, HStack, Input, Text } from '@chakra-ui/react'
import { useFormContext } from 'react-hook-form'

const TunerStep = () => {
  const {
    setValue,
    register,
    formState: { errors },
    watch,
  } = useFormContext();

  const watchTuner = watch("tuner");

  return (
    <Box p="10" pb="0">
      <Heading as="h3" mb="2" size="md">
        Number of tuners
      </Heading>
      <Text mb="5" fontSize="sm">
        Number of simultaneous connections that can be established to the
        playlist provider.
      </Text>
      <Alert status="warning" mb="10" variant="left-accent">
        <AlertIcon />
        <Text fontSize="xs">
          Available for Plex, Emby (HDHR), M3U (with stream buffer enabled)
        </Text>
      </Alert>
      <HStack>
        <Button
          colorScheme="teal"
          disabled={parseInt(watchTuner, 10) === 1}
          size="lg"
          onClick={() =>
            setValue(
              "tuner",
              parseInt(watchTuner, 10) > 1 ? parseInt(watchTuner, 10) - 1 : 1
            )
          }
        >
          -
        </Button>
        <FormControl isInvalid={errors.tuner}>
          <Input
            size="lg"
            type="number"
            min="1"
            {...register("tuner", {
              required: "This is required",
              min: 1,
            })}
          />
        </FormControl>
        <Button
          colorScheme="teal"
          size="lg"
          onClick={() => setValue("tuner", parseInt(watchTuner, 10) + 1)}
        >
          +
        </Button>
      </HStack>
      <FormErrorMessage>
        {errors.tuner && errors.tuner.message}
      </FormErrorMessage>
    </Box>
  );
};

export default TunerStep;
