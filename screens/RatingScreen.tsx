import React, { useEffect, useState } from 'react';
import { View, StyleSheet, TouchableOpacity, SafeAreaView } from 'react-native';
import { Avatar, Button, Text, Card, Title, Paragraph, Appbar } from 'react-native-paper';
import { Rating } from 'react-native-ratings';
import { useNavigation } from '@react-navigation/native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';

type RootStackParamList = {
  NotificationScreen: undefined;
  RatingScreen: undefined;
  RatingSucceedScreen: undefined;
};

type PostOwnerProfile = { 
  post_id: string;
  name: string;
  post_rating: number;
}

const defaultPostOwnerProfile: PostOwnerProfile = {
  post_id: '2',
  name: 'Jimmy Ho',
  post_rating: 0,
};

type NotificationsProps = NativeStackScreenProps<RootStackParamList, 'NotificationScreen'>;

const RatingScreen: React.FC<NotificationsProps> = ({ navigation }: NotificationsProps) => {

  const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic });

  const [postOwnerProfile, setPostOwnerProfile] = useState<PostOwnerProfile>(defaultPostOwnerProfile);

  useEffect(() => {                                                                                     //fill this in to get db info 
    const fetchPostOwnerProfile = async () => {                                                         //Set rating default to 0, then edit it after receiving user rating in this screen
      try {
        const response = await fetch('URL_TO_YOUR_BACKEND/postOwnerProfile_endpoint');
        const json = await response.json();
        setPostOwnerProfile(json); // Adjust this depending on the structure of your JSON
      } catch (error) {
        // console.error(error); // uncomment this after finish frontend developing
      }
    };
    fetchPostOwnerProfile();
  }, []);

  const use_navigation = useNavigation(); //for Appbar.BackAction

  const [rating, setRating] = useState<number>(0);

  // update postOwnerProfile.post_rating to the received rating value
  function updateProfile(currentInfo: PostOwnerProfile, updates: Partial<PostOwnerProfile>): PostOwnerProfile {
    return {
      ...currentInfo,
      ...updates,
    };
  }

  const submitRating = () => {                                           
    setPostOwnerProfile(prev => updateProfile(prev, { post_rating: rating }))
    console.log('Rating submitted: ', rating);
    navigation.navigate('RatingSucceedScreen');
  };

  useEffect(() => {
    // This will log the updated rating only when postOwnerProfile.post_rating changes
    console.log('Updated Rating in postOwnerProfile: ', postOwnerProfile.post_rating);
  }, [postOwnerProfile.post_rating]);

  // Need to submit the edited postOwnerProfile (with updated postOwnerProfile.post_rating) to backend; MAKE SURE the api calling is inside useEffect so that the new value could be fetched

  return (
    <SafeAreaView  style={styles.container}> 
      <View style={styles.headerContainer}>
        <Appbar.BackAction style={styles.backAction} onPress={() => use_navigation.goBack()} />
        <Text style={styles.header}>GiveGetGo</Text>
        <View style={styles.backActionPlaceholder} />
      </View>
      <Card style={styles.card}>
        <Card.Content>
          <Text style={styles.question}>
            How would you rate your match with <Paragraph style={styles.boldQuestion}>{postOwnerProfile.name}</Paragraph>?
          </Text>
          <Rating
            type="custom"
            ratingCount={5}
            imageSize={40} // Smaller size for the stars
            fractions={0} // Set the granularity of the rating; 1 means full stars
            startingValue={rating} // The initial rating value
            tintColor="#f6f6f6" // Background color for the rating component
            ratingBackgroundColor="#c8c7c8" // Color behind the non-selected rating icons
            ratingColor="orange" // Color of the selected rating icons
            showRating
            onFinishRating={(rating: number) => setRating(rating)}
            style={styles.rating}
          />
          <Button style={styles.button} mode="contained" onPress={submitRating}>
            Submit
          </Button>
        </Card.Content>
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
    right: 0,
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
    height: 48,
  },
  backAction: {
    marginLeft: 0 //This means the relative margin, comparing to the container (?)
  },
  card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 15, // Add padding inside the card
    height: 250,
  },
  button: {
    position: 'absolute', 
    left: 110,
    right: 110, //position, left, right together controls the button's length and horizontal location
    bottom:-35,
    alignSelf: 'center', 
  },
  question: {
    fontSize: 17,
    textAlign: 'center',
    marginTop: -10,
    marginBottom: -13,
    padding: 10,
  },
  boldQuestion: {
    fontSize: 17,
    textAlign: 'center',
    marginBottom: -13,
    fontWeight: 'bold',
    padding: 10,
  },
  rating: {
    paddingVertical: 10,
  },
});

export default RatingScreen;
