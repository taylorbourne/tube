import { ArrowForwardIcon } from '@chakra-ui/icons'
import { Box, Button, Center, Container, Heading, SlideFade, Text } from '@chakra-ui/react'
import { Step, Steps, useSteps } from 'chakra-ui-steps'
import React, { useEffect, useState } from 'react'
import { FormProvider, useForm } from 'react-hook-form'

import WizardProvider from '../../context/Wizard'
import EPGStep from './EPGStep'
import PlaylistStep from './PlaylistStep'
import StepButtons from './StepButtons'
import TunerStep from './TunerStep'
import XMLTVStep from './XMLTVStep'

const Base = ({ startWizard }: { startWizard: () => void }) => (
  <>
    <Heading as="h1" mb="5" textAlign="center">
      Welcome to XTV
    </Heading>
    <Text mb="10" textAlign="center">
      Now that your server is up and running let's get things configured. Click
      the button below to start the wizard and configure your first playlist.
      You can add more playlists later.
    </Text>
    <Center>
      <Button
        rightIcon={<ArrowForwardIcon />}
        colorScheme="teal"
        onClick={() => startWizard()}
        variant="outline"
      >
        Get Started
      </Button>
    </Center>
  </>
);

const conditionalXMLTVStep = {
  label: "XMLTV",
  content: <XMLTVStep />,
  description: "Initial xmltv setup",
};

const baseSteps = [
  {
    description: "Provider setup",
    content: <TunerStep />,
    label: "Tuners",
  },
  { label: "EPG", content: <EPGStep />, description: "Configure EPG" },
  {
    label: "Playlist",
    content: <PlaylistStep />,
    description: "Initial playlist setup",
  },
  conditionalXMLTVStep,
];

const SetupWizard = () => {
  const methods = useForm({
    defaultValues: {
      tuner: 1,
      epgSource: true,
    },
    mode: "onChange",
  });

  const { formState, watch } = methods;
  const [isStarted, setIsStarted] = useState(false);
  const [steps, setSteps] = useState(baseSteps);
  const { nextStep, prevStep, setStep, reset, activeStep } = useSteps({
    initialStep: 0,
  });

  const watchEPGSource = watch("epgSource");

  useEffect(() => {
    const stepsSet = new Set(steps);
    if (watchEPGSource) {
      if (!steps.includes(conditionalXMLTVStep)) {
        steps.push(conditionalXMLTVStep);
      }
    } else {
      setSteps((prevSteps) =>
        prevSteps.filter((step) => step.label !== "XMLTV")
      );
    }
  }, [watchEPGSource]);

  const onSubmit = methods.handleSubmit((data) => {
    nextStep();
    console.log("data", data);
  });

  return (
    <WizardProvider>
      <FormProvider {...methods}>
        <Container
          alignItems="center"
          justifyContent="center"
          centerContent
          h="100vh"
          maxW={{ base: "100%", md: "container.md", lg: "container.lg" }}
          p={{ base: 0 }}
        >
          <form onSubmit={onSubmit}>
            <Box
              bgColor="gray.900"
              w={{ base: "100vw", md: "auto" }}
              h={{ base: "100vh", md: "auto" }}
              borderRadius="md"
              justifyContent="center"
              display="flex"
              flexDirection="column"
              overflow="hidden"
            >
              <Box p={{ base: "5", lg: "10" }}>
                {!isStarted && (
                  <SlideFade in={!isStarted} offsetY="20px">
                    <Base startWizard={() => setIsStarted(true)} />
                  </SlideFade>
                )}
                {isStarted && (
                  <SlideFade in={isStarted} offsetY="20px">
                    <Box>
                      <Steps
                        activeStep={activeStep}
                        colorScheme="teal"
                        size="sm"
                      >
                        {steps.map(({ label, content, description }, index) => (
                          <Step
                            description={description}
                            label={label}
                            key={label}
                          >
                            <SlideFade in={activeStep === index} offsetY="20px">
                              <Box px={{ base: "0", lg: "20" }} py="10">
                                {content}
                              </Box>
                            </SlideFade>
                          </Step>
                        ))}
                      </Steps>
                    </Box>
                  </SlideFade>
                )}
              </Box>
              {isStarted &&
                (activeStep === steps.length ? (
                  <>doingstuff</>
                ) : (
                  <StepButtons
                    {...{ activeStep, nextStep, prevStep }}
                    isFirstStep={activeStep === 0}
                    isFinalStep={activeStep === steps.length - 1}
                    isFormValid={formState.isValid}
                    nextDisabled={activeStep === steps.length - 1}
                    prevDisabled={activeStep === 0}
                  />
                ))}
            </Box>
          </form>
        </Container>
      </FormProvider>
    </WizardProvider>
  );
};
export default SetupWizard;
