import { Alert, AlertIcon, Box, FormControl, FormHelperText, Heading, Input, Text } from '@chakra-ui/react'
import { useFormContext } from 'react-hook-form'

const XMLTVStep = () => {
  const {
    register,
    formState: { errors },
  } = useFormContext();
  return (
    <Box p={{ base: 5, md: 10 }} pl={{ base: 0, md: 10 }} pb="0">
      <Heading as="h3" mb="2" size="md">
        XMLTV URL
      </Heading>
      <Text mb="5" fontSize="sm">
        Since you have enabled the XEPG feature, XTV will manage your channel
        and XMLTV data. Enter your XMLTV URL below to begin importing.
      </Text>
      <FormControl id="xmltv-url" isInvalid={errors.xmltv} mb="10">
        <Input {...register("xmltv", { required: true })} type="text" />
        <FormHelperText>
          Enter a valid local or remote path to your XMLTV file
        </FormHelperText>
      </FormControl>
      <Alert status="info" variant="left-accent">
        <AlertIcon />
        <Text fontSize="xs">
          We'll just add one XMLTV file for now, but you can add as many as you
          would like later.
        </Text>
      </Alert>
    </Box>
  );
};

export default XMLTVStep;
