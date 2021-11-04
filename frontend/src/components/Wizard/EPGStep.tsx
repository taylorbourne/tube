import { CheckCircleIcon } from '@chakra-ui/icons'
import {
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
  Box,
  Heading,
  List,
  ListIcon,
  ListItem,
  Switch,
  Text,
} from '@chakra-ui/react'
import { useFormContext } from 'react-hook-form'

const EPGStep = () => {
  const { register, watch } = useFormContext();

  const watchEpgSource = watch("epgSource");

  return (
    <Box p="10" pb="0">
      <Heading as="h3" mb="2" size="md">
        EPG Setup
      </Heading>
      <Text mb="5" fontSize="sm">
        Choose weather or not to enable XEPG. XEPG enables users to merge data
        from multiple XMLTV sources into a single EPG, with full control over
        channel numbering and naming.
      </Text>

      <Box
        alignItems="center"
        bgColor="rgba(0, 0, 0, 0.25)"
        borderRadius="md"
        display="flex"
        flexDirection="row"
        p="3"
      >
        <Switch
          colorScheme="teal"
          defaultChecked={watchEpgSource}
          size="lg"
          {...register("epgSource")}
        />
        <Text
          ml="5"
          color={watchEpgSource ? "teal.400" : "red.300"}
          fontSize="sm"
          fontWeight="bold"
        >
          XEPG {watchEpgSource ? "Enabled" : "Disabled"}
        </Text>
      </Box>

      {watchEpgSource && (
        <Alert status="info" mt="10" variant="left-accent">
          <AlertIcon />
          <AlertDescription>
            <AlertTitle fontSize="xs">
              With XEPG enabled, you will have access to the following features:
            </AlertTitle>
            <List fontSize="xs" size="sm" spacing={1}>
              <ListItem alignItems="center" display="flex">
                <ListIcon as={CheckCircleIcon} color="green.500" />
                Use of one or more EPG sources
              </ListItem>
              <ListItem>
                <ListIcon as={CheckCircleIcon} color="green.500" />
                Channel management
              </ListItem>
              <ListItem>
                <ListIcon as={CheckCircleIcon} color="green.500" />
                M3U/XMLTV external link
              </ListItem>
            </List>
          </AlertDescription>
        </Alert>
      )}

      {!watchEpgSource && (
        <Alert fontSize="xs" status="warning" mt="10" variant="left-accent">
          <AlertIcon />
          <AlertDescription>
            <AlertTitle>XEPG is disabled</AlertTitle>
            With XEPG disabled, you will have to rely on EPG data from your own
            source such as Plex, Emby, or Jellyfin.
          </AlertDescription>
        </Alert>
      )}
    </Box>
  );
};

export default EPGStep;
