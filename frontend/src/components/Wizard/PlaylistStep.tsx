import { Alert, AlertIcon, Box, FormControl, FormHelperText, Heading, Input, Text } from '@chakra-ui/react'
import { useFormContext } from 'react-hook-form'

const PlaylistStep = () => {
  const {
    register,
    formState: { errors },
  } = useFormContext();
  return (
    <Box p={{ base: 5, md: 10 }} pb="0">
      <Heading as="h3" mb="2" size="md">
        Playlist URL
      </Heading>
      <Text mb="5" fontSize="sm">
        XTV can parse your M3U playlist from any remote URL or local path. Be
        sure to include the proper protocal (http or https, for example) when
        using a remote source.
      </Text>
      <FormControl isInvalid={errors.m3u} mb="10">
        <Input {...register("m3u", { required: true })} type="text" />
        <FormHelperText>
          Enter a valid local or remote path to your playlist
        </FormHelperText>
      </FormControl>
      <Alert status="info" variant="left-accent">
        <AlertIcon />
        <Text fontSize="xs">
          We'll just add one playlist for now, but you can add as many as you
          would like later.
        </Text>
      </Alert>
    </Box>
  );
};

export default PlaylistStep;
