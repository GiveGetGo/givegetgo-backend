import React from 'react';
import { View, StyleSheet, SafeAreaView } from 'react-native';
import { Button, Text, Card, Title, Paragraph } from 'react-native-paper';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { useRoute } from '@react-navigation/native';
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';

type RootStackParamList = {
  HomeScreen: undefined;
  PostRequestSucceedScreen: undefined; 
};

type HomeScreenProps = NativeStackScreenProps<RootStackParamList, 'HomeScreen'>;

const PostRequestSucceedScreen: React.FC<HomeScreenProps> = ({ navigation }: HomeScreenProps) => {

  const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic });

  const route = useRoute();
  const { postId } = route.params as { postId: string };
  const { name } = route.params as { name: string };

  console.log("postID parsed to PostRequestSucceedScreen: ", postId) 
  console.log("postOwnerName parsed to PostRequestSucceedScreen: ", name) 

  const GetContact = () => {
    navigation.navigate('HomeScreen');
  };

  return (
    <SafeAreaView  style={styles.container}> 
      <View style={styles.headerContainer}>
        <View style={styles.backActionPlaceholder} />
        <Text style={styles.header}>GiveGetGo</Text>
        <View style={styles.backActionPlaceholder} />
      </View>
      <Card style={styles.card}>
        <Card.Content>
          <Title style={styles.title}>Congratulations!</Title>
          <Paragraph style={styles.paragraph_firstline}>Your request has been submitted. </Paragraph>
          <Paragraph style={styles.paragraph}>         
            You will be able to contact <Paragraph style={styles.boldText}>{name}</Paragraph> once a match has been established.
          </Paragraph>
        </Card.Content>
        <Card.Actions style={styles.cardActions}>
          <Button style={styles.button} mode="contained" onPress={GetContact}>
            Home
          </Button>
        </Card.Actions>
      </Card>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,                                
    marginTop: 50,
    justifyContent: 'center',
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    justifyContent: 'space-between', // Distributes items evenly horizontally
    paddingLeft: 10, 
    paddingRight: 10, 
    position: 'absolute', // So that while setting card to the vertical middle, it still stays at the same place
    top: 0, 
    left: 0,
    right: 2,
    zIndex: 1, // Ensure the headerContainer is above the card
  },
  header: {
    fontSize: 22, // Increase the font size
    fontWeight: '600', // Make the font weight bold
    fontFamily: 'Montserrat_700Bold_Italic',
    textAlign: 'center', // Center the text
    color: '#444444', // Dark gray color
  },
  backActionPlaceholder: {
    width: 48, // This should match the width of the Appbar.BackAction for balance
    height: 52,
  },
  card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 15, // Add padding inside the card
  },
  title: {
    textAlign: 'center',
    fontWeight: 'bold',
    // marginVertical: 3,
    marginBottom: 5,
    marginTop: -10,
  },
  paragraph: {
    textAlign: 'center',
    fontSize: 16,
    marginBottom: 12,
  },
  paragraph_firstline: {
    textAlign: 'center',
    fontSize: 16,
    marginTop: -3,
    marginBottom: 0,
  },
  boldText: {
    fontWeight: 'bold',
    textAlign: 'center',
    fontSize: 16,
    marginBottom: 12,
  },
  button: {
    position: 'absolute', 
    left: 110,
    right: 110, //position, left, right together controls the button's length and horizontal location
    alignSelf: 'center', 
  },
  cardActions: {
    justifyContent: 'center', 
    alignItems: 'center',
    padding: 20,
  },
});

export default PostRequestSucceedScreen;
