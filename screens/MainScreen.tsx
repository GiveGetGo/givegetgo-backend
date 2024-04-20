import React, { useState, useEffect } from 'react';
import { SafeAreaView, StyleSheet, View, FlatList, Dimensions, TouchableOpacity, TouchableWithoutFeedback, Keyboard } from 'react-native';
import { BottomNavigation, Text, Searchbar, Card, Title, Paragraph, Avatar, IconButton, Button, List, Divider, TextInput } from 'react-native-paper';
import { NavigationContainer, CommonActions } from '@react-navigation/native';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { NativeStackScreenProps, createNativeStackNavigator } from '@react-navigation/native-stack';
import { useRoute, RouteProp } from '@react-navigation/native';
import { useSelector } from 'react-redux';
import SeeRequestScreen from './SeeRequestScreen'; 
import RatingScreen from './RatingScreen'; 
import ProfileScreen from './ProfileScreen'; 
import RequestSucceedScreen from './RequestSucceedScreen'; 
import GiveOutContactScreen from './GiveOutContactScreen'; 
import RatingSucceedScreen from './RatingSucceedScreen'; 
import NotificationStackProfileScreen from './NotificationStackProfileScreen'; 
import PostDetailsScreen from './PostDetailsScreen'; 
import PostRequestInfoScreen from './PostRequestInfoScreen'; 
import PostRequestSucceedScreen from './PostRequestSucceedScreen'; 
import AvatarPickerScreen from './AvatarPickerScreen'; 

type RootStackParamList = {
    MainScreen: undefined;
    ProfileScreen: undefined;
    FilterScreen: undefined;
    NotificationScreen: undefined;
    SeeRequestScreen: undefined;
    RatingScreen: undefined;
    HomeScreen: undefined;
    ProfileSreen: undefined;
    othersProfileScreen: undefined;
    RequestSucceedScreen: undefined;
    GiveOutContactScreen: undefined;
    RatingSucceedScreen: undefined;
    PostDetailsScreen: { postId: string };
    PostRequestInfoScreen: undefined; 
    PostRequestSucceedScreen: undefined;
    SettingsScreen: {newAvatarUri: string};
    AvatarPickerScreen: undefined;
  };

type Post = {
    id: string;
    title?: string; // Optional, if your posts have titles
    description?: string; // Optional, if your posts have descriptions
    imageUri?: string;
    // Include other properties for your posts, like images, etc.
  };

type Notification = {
    id: string;
    title?: string;
    description?: string;
    time?: string; // This should be a date string
};

type UserInfo = {
    name: string;
    email: string;
    bio: string;
    profilePicture: string; // This should be a URI string
    major: string;
    classYear: number;
  };

type MainScreenProps = {
    // Define props here if any
};

type HomeScreenProps = NativeStackScreenProps<RootStackParamList, 'HomeScreen'>;

type MakingPostsProps = {
    // Define props here if any
};

type NotificationsProps = NativeStackScreenProps<RootStackParamList, 'NotificationScreen'>;

type SettingsScreenProps = NativeStackScreenProps<RootStackParamList, 'SettingsScreen'>;

// dummy data for posts //should be loaded from the database
const posts: Post[] = [
    { id: '1', title: 'Post 11111111111', description: 'Description 11111111111111111111111111'},
    { id: '2', title: 'Post 22222222222', description: 'Description 2222'},
    { id: '3', title: 'Post 3', description: 'Description 3'},
    // ...more posts
  ];

// dummy data for notifications //should be loaded from the database
const notifications_data: Notification[] = [
    {
      id: '1',
      title: 'Request',
      time: '1 hour ago',
      description: 'You have an matching request from xxxxxxxxxx for your xxxxxxx ',
    },
    {
      id: '2',
      title: 'Match',
      description: 'Match succeeded with xxxxx  for xxxxxxxx. Click in to rate this match!',
      time: '2 hour ago',
    },
  ];

// dummy data for settings is written locally

const Tab = createBottomTabNavigator();

// Build Home's Stack
const HomeStack = createNativeStackNavigator();
function HomeStackScreen() {
  return (
    <HomeStack.Navigator>
      <HomeStack.Screen name="HomeScreen" component={HomeScreen} options={{ headerShown: false }}/>
      <HomeStack.Screen name="ProfileScreen" component={ProfileScreen} options={{ headerShown: false }}/>
      <HomeStack.Screen name="PostDetailsScreen" component={PostDetailsScreen} options={{ headerShown: false }}/>
      <HomeStack.Screen name="PostRequestInfoScreen" component={PostRequestInfoScreen} options={{ headerShown: false }}/>
      <HomeStack.Screen name="PostRequestSucceedScreen" component={PostRequestSucceedScreen} options={{ headerShown: false }}/>
    </HomeStack.Navigator>
  );
}

// Build Notification's Stack
const NotificationsStack = createNativeStackNavigator();
function NotificationsStackScreen() {
  return (
    <NotificationsStack.Navigator initialRouteName="NotificationScreen">
      <NotificationsStack.Screen name="NotificationScreen" component={NotificationScreen} options={{ headerShown: false }}/>
      <NotificationsStack.Screen name="SeeRequestScreen" component={SeeRequestScreen} options={{ headerShown: false }}/>
      <NotificationsStack.Screen name="NotificationStackProfileScreen" component={NotificationStackProfileScreen} options={{ headerShown: false }}/>
      <NotificationsStack.Screen name="RatingScreen" component={RatingScreen} options={{ headerShown: false }}/>
      <NotificationsStack.Screen name="RequestSucceedScreen" component={RequestSucceedScreen} options={{ headerShown: false }}/>
      <NotificationsStack.Screen name="GiveOutContactScreen" component={GiveOutContactScreen} options={{ headerShown: false }}/>
      <NotificationsStack.Screen name="RatingSucceedScreen" component={RatingSucceedScreen} options={{ headerShown: false }}/>
    </NotificationsStack.Navigator>
  );
}

// Build Settings's Stack
const SettingsStack = createNativeStackNavigator();
function SettingsStackScreen() {
  return (
    <SettingsStack.Navigator initialRouteName="SettingsScreen">
      <SettingsStack.Screen name="SettingsScreen" component={SettingsScreen} options={{ headerShown: false }}/>
      <SettingsStack.Screen name="AvatarPickerScreen" component={AvatarPickerScreen} options={{ headerShown: false }}/>
    </SettingsStack.Navigator>
  );
}

const MainScreen: React.FC<MainScreenProps> = () => {
    return (
        <NavigationContainer independent={true}>
            <Tab.Navigator
                screenOptions={{
                headerShown: false,
                }}
                tabBar={({ navigation, state, descriptors, insets }) => (
                    <BottomNavigation.Bar
                        navigationState={state}
                    safeAreaInsets={insets}
                        onTabPress={({ route, preventDefault }) => {
                        const event = navigation.emit({
                            type: 'tabPress',
                            target: route.key,
                            canPreventDefault: true,
                        });
        
                        if (event.defaultPrevented) {
                            preventDefault();
                        } else {
                        navigation.dispatch({
                            ...CommonActions.navigate(route.name, route.params),
                            target: state.key,
                            });
                        }
                        }}
                        renderIcon={({ route, focused, color }) => {
                        const { options } = descriptors[route.key];
                        if (options.tabBarIcon) {
                            return options.tabBarIcon({ focused, color, size: 24 });
                        }
        
                        return null;
                        }}
                        getLabelText={({ route }) => {
                        const { options } = descriptors[route.key];
                        const label =
                            options.tabBarLabel !== undefined
                            ? options.tabBarLabel
                            : options.title !== undefined
                            ? options.title
                            : route.title;
        
                        return label;
                        }}
                    />
                    )}
            >
                <Tab.Screen
                name="Home"
                component={HomeStackScreen}
                options={{
                    tabBarLabel: 'Home',
                    tabBarIcon: ({ color, size }) => {
                    return <Icon name="home" size={size} color={color} />;
                    },
                }}
                />
                <Tab.Screen
                name="Post"
                component={PostScreen}
                options={{
                    tabBarLabel: 'Post',
                    tabBarIcon: ({ color, size }) => {
                    return <Icon name="message-plus" size={size} color={color} />;
                    },
                }}
                />
                <Tab.Screen
                name="Notifications"
                component={NotificationsStackScreen}
                options={{
                    tabBarLabel: 'Notifications',
                    tabBarIcon: ({ color, size }) => {
                    return <Icon name="bell-ring" size={size} color={color} />;
                    },
                }}
                />
                <Tab.Screen
                name="Settings"
                component={SettingsStackScreen}
                options={{
                    tabBarLabel: 'Settings',
                    tabBarIcon: ({ color, size }) => {
                    return <Icon name="account-settings" size={size} color={color} />;
                    },
                }}
                />
            </Tab.Navigator>
        </NavigationContainer>
  );
};

const HomeScreen: React.FC<HomeScreenProps> = ({ navigation }: HomeScreenProps) => {
    const [searchQuery, setSearchQuery] = React.useState('');
    const onChangeSearch = (query: string) => setSearchQuery(query);
    const searchResults = ['Post 1', 'Post 3', 'Post 5']; // Example static data
    // const navigation = useNavigation<ScreenNavigationProp>();

    const handleProfilePress = () => {
        navigation.navigate('ProfileScreen'); 
    };

    const filteredPosts = posts.filter(post =>
        post.title?.toLowerCase().includes(searchQuery.toLowerCase())
      );
    
    const navigateToPostDetails = (postId: string) => {
        navigation.navigate('PostDetailsScreen', {postId: postId,});
    };

    const numColumns = 2;
    const { width } = Dimensions.get('window');
    const cardWidth = width / numColumns - 16;                     
    const renderItem = ({ item }: { item: Post }) => (
        <Card style={[styles.card, { width: cardWidth }]} onPress={() => navigateToPostDetails(item.id)} >
            {/* {item.imageUri && <Card.Cover source={{ uri: item.imageUri }} />}   */}
            <Card.Content>
            <Title>{item.title}</Title>
            <Paragraph>{item.description}</Paragraph>
            </Card.Content>
        </Card>
    );

    // get selectedAvatarUri from Redux
    const selectedAvatarUri = useSelector((state: { avatar: { selectedAvatarUri: string } }) => state.avatar.selectedAvatarUri); 
    console.log("selectedAvatarUri from redux:", selectedAvatarUri) 


    // require() could not directly take the url
    interface ImageMap {
        [key: string]: NodeRequire;
    }
    const imageMap: { [key: string]: NodeRequire } = {
        '../assets/avatars/avatar1.png': require('../assets/avatars/avatar1.png'),
        '../assets/avatars/avatar2.png': require('../assets/avatars/avatar2.png'),
        '../assets/avatars/avatar3.png': require('../assets/avatars/avatar3.png'),
        '../assets/avatars/avatar4.png': require('../assets/avatars/avatar4.png'),
        '../assets/avatars/avatar5.png': require('../assets/avatars/avatar5.png'),
        '../assets/avatars/avatar6.png': require('../assets/avatars/avatar6.png'),
        '../assets/avatars/avatar7.png': require('../assets/avatars/avatar7.png'),
        '../assets/avatars/avatar8.png': require('../assets/avatars/avatar8.png'),
        '../assets/avatars/avatar9.png': require('../assets/avatars/avatar9.png'),
    };
    const getProfilePictureSource = (uri: string) => {
        return imageMap[uri] || require(`./profile_icon.jpg`);
    };


    return (
        <SafeAreaView  style={styles.container}>
            <View style={styles.headerContainer}>
                <TouchableOpacity onPress={handleProfilePress}>
                    <Avatar.Image size={40} source={getProfilePictureSource(selectedAvatarUri)} />
                </TouchableOpacity>
                {/* <TouchableOpacity onPress={handleFilterPress}>
                    <IconButton
                        icon="filter"
                        size={24}
                        onPress={handleFilterPress}
                        style={styles.filterButton}
                    />
                </TouchableOpacity> */}
                <Text style={styles.header}>GiveGetGo</Text>
            </View>
            <Searchbar
                placeholder="Search"
                onChangeText={onChangeSearch}
                value={searchQuery}
                style={styles.searchbar}
            />
            <FlatList
                data={filteredPosts}
                renderItem={renderItem}
                // renderItem={renderItem}
                keyExtractor={(item) => item.id}
                numColumns={numColumns}
                columnWrapperStyle={styles.column}
            />
        </SafeAreaView >
    );
};

const PostScreen: React.FC<MakingPostsProps> = () => {
    // const handleImageUploadPress = () => {
    // // TODO: Implement your image upload logic here (might be using "launchImageLibrary" but i could not get it work; stacloverflow shows this error happens only on ios aka expo go)
    //     const options = {
    //         mediaType: 'photo' as const,
    //         quality: 1,
    //     };
    //   // Launch the image picker
    //   launchImageLibrary(option, (response) => {
    //     if (response.didCancel) {
    //       console.log('User cancelled image picker');
    //     } else if (response.error) {
    //       console.log('ImagePicker Error: ', response.error);
    //     } else {
    //       // The response object contains various information about the selected image
    //       const source = { uri: response.uri };
    
    //       // TODO: Implement your upload logic here
    //       // For example, you could upload the image to a server using a POST request
    //       console.log('Selected image: ', source);
    //     }
    //   });
    // };

  const handlePostPress = () => {
    // TODO: Implement your submit post logic here
    console.log('Post submitted with Title:', title, 'and Description:', description);
  };

  // Placeholder for your image icon, replace 'require' with the actual path to your icon
  const imageUploadIcon = require('./image_icon.jpg');

  const [title, onChangeTitle] = React.useState<string>('');
  const [description, onChangeDescription] = React.useState<string>('');
    return (
        <TouchableWithoutFeedback onPress={Keyboard.dismiss} accessible={false}>
            <SafeAreaView style={styles.container}>
                <View style={styles.headerContainer}>
                    <Text style={styles.header}>GiveGetGo</Text>
                </View>
                <Card>
                    <Card.Title title="Create a Post" right={(props) => <Button {...props} onPress={handlePostPress}>Post</Button>} />
                    <Card.Content>
                        <TextInput
                            editable
                            style={styles.input}
                            onChangeText={onChangeTitle}
                            placeholder="Add a title..."
                            value={title}
                            returnKeyType="done"
                        />
                        <TextInput
                            editable
                            multiline
                            numberOfLines={4}
                            maxLength={150}
                            onChangeText={onChangeDescription}
                            placeholder="Add a description..."
                            value={description}
                            returnKeyType="done"
                        />
                    </Card.Content>
                    {/* <Card.Actions> */}
                        {/* <IconButton
                            icon="image"
                            size={24}
                            onPress={handleImageUploadPress}
                        /> */}
                    {/* </Card.Actions> */}
                </Card>
            </SafeAreaView>
        </TouchableWithoutFeedback>
    );
};

const NotificationScreen: React.FC<NotificationsProps> = ({ navigation }: NotificationsProps) => {
    const [notifications, setNotifications] = React.useState<Notification[]>(notifications_data);                   

    const deleteNotification = (id: string) => {
        setNotifications(currentNotifications => currentNotifications.filter(notification => notification.id !== id));

        // build api here to remove the corresponding notification in the backend/database
    };

    const renderItem = ({ item }: { item: Notification })  => (
        <Card
          style={styles.notifications_card}
          onPress={() => {
            // Navigate based on item.title
            if (item.title === 'Request') {
                navigation.navigate('SeeRequestScreen')
            } else if (item.title === 'Match') {
                navigation.navigate('RatingScreen')
            }
          }}
        >
        <Card.Content>
            <Title style={styles.notifications_title}>{item.title}</Title>
            <Paragraph style={styles.notifications_time}>{item.time}</Paragraph>
            <Paragraph style={styles.notifications_description}>{item.description}</Paragraph>
        </Card.Content>
        <Card.Actions>
            <Paragraph style={{ flex: 1 }}>...</Paragraph>
            <IconButton
                icon="trash-can" 
                size={18} 
                style={styles.deleteIcon}
                onPress={() => deleteNotification(item.id)}
            />
        </Card.Actions>
        </Card>
      );

    return (
        <SafeAreaView  style={styles.container}>
            <View style={styles.headerContainer}>
                <Text style={styles.header}>GiveGetGo</Text>
            </View>
            <FlatList
                data={notifications}
                renderItem={renderItem}
                keyExtractor={(item) => item.id}
                // removeClippedSubviews={true} // If you want to improve performance on large lists, uncomment the line
            />
        </SafeAreaView >
    );
};

const SettingsScreen: React.FC<SettingsScreenProps> = ({ navigation }: SettingsScreenProps) => {  

    let currentUserInfo: UserInfo = { //should be loaded from the database
        name: 'Gilbert Hsu',
        email: 'xxx@gmail.com',
        bio: 'I am Gilbert the magician.',
        profilePicture: './profile_icon.jpg', // Placeholder for profile picture URI
        major: 'Computer Science',
        classYear: 2024,
      };

    const [userInfo, setUserInfo] = React.useState<UserInfo>(currentUserInfo);

    function updateUserInfo(currentInfo: UserInfo, updates: Partial<UserInfo>): UserInfo {
        // Return a new object that combines the current info with the updates
        // The spread operator (...) is used to first copy all current properties
        // Then, properties in `updates` will overwrite or add to those in `currentInfo`
        return {
          ...currentInfo,
          ...updates,
        };
      }

    const [isEditing_name, setIsEditing_name] = useState(false);
    const [name, setName] = useState<string>(userInfo.name);
    const [isEditing_email, setIsEditing_email] = useState(false);
    const [email, setEmail] = useState<string>(userInfo.email);
    const [isEditing_bio, setIsEditing_bio] = useState(false);
    const [bio, setBio] = useState<string>(userInfo.bio);
    const [isEditing_pic, setIsEditing_pic] = useState(false);
    const [profilePicture, setProfilePicture] = useState<string>(userInfo.profilePicture);

    const handleEditPress_name = () => {
        setIsEditing_name(true);
    };

    const handleSavePress_name = () => {
        setUserInfo(prev => updateUserInfo(prev, { name: name })); //updated userInfo should also be saved to the database
        setIsEditing_name(false);
    };

    const handleEditPress_email = () => {
        setIsEditing_email(true);
    };

    const handleSavePress_email = () => {
        setUserInfo(prev => updateUserInfo(prev, { email: email }));
        setIsEditing_email(false);
    };

    const handleEditPress_bio = () => {
        setIsEditing_bio(true);
    };

    const handleSavePress_bio = () => {
        setUserInfo(prev => updateUserInfo(prev, { bio: bio }));
        setIsEditing_bio(false);
    };

    const handleEditPress_pic = () => {
        setIsEditing_pic(true);
        navigation.navigate('AvatarPickerScreen')  
    };
    
    const handleLogOut = () => {
        console.log('Log out');
    };

    // Get newAvatarUri from AvatarPickerScreen
    type SettingsScreenRouteProp = RouteProp<RootStackParamList, 'SettingsScreen'>;
    const route = useRoute<SettingsScreenRouteProp>();
    const newAvatarUri = route.params?.newAvatarUri;
    console.log('New Avatar URI:', newAvatarUri);

    // Update profile picture when newAvatarUri changes
    useEffect(() => {
        if (newAvatarUri) {
            console.log("newAvatarUri detected")
            setUserInfo(prevUserInfo => ({
                ...prevUserInfo,
                profilePicture: newAvatarUri
            }));
        }
    }, [newAvatarUri]);


    // require() could not directly take the url
    interface ImageMap {
        [key: string]: NodeRequire;
    }
    const imageMap: { [key: string]: NodeRequire } = {
        '../assets/avatars/avatar1.png': require('../assets/avatars/avatar1.png'),
        '../assets/avatars/avatar2.png': require('../assets/avatars/avatar2.png'),
        '../assets/avatars/avatar3.png': require('../assets/avatars/avatar3.png'),
        '../assets/avatars/avatar4.png': require('../assets/avatars/avatar4.png'),
        '../assets/avatars/avatar5.png': require('../assets/avatars/avatar5.png'),
        '../assets/avatars/avatar6.png': require('../assets/avatars/avatar6.png'),
        '../assets/avatars/avatar7.png': require('../assets/avatars/avatar7.png'),
        '../assets/avatars/avatar8.png': require('../assets/avatars/avatar8.png'),
        '../assets/avatars/avatar9.png': require('../assets/avatars/avatar9.png'),
    };
    const getProfilePictureSource = (uri: string) => {
        return imageMap[uri] || require(`./profile_icon.jpg`);
    };


    // get animal names for description
    interface AnimalMapping {
        [key: string]: string;
    }
    const animalMapping: AnimalMapping = {
        'avatar1.png': 'Fox',
        'avatar2.png': 'Owl',
        'avatar3.png': 'Dog',
        'avatar4.png': 'Seal',
        'avatar5.png': 'Duck',
        'avatar6.png': 'Cat',
        'avatar7.png': 'Goat',
        'avatar8.png': 'Lion',
        'avatar9.png': 'Penguin',
    };
    

    function mapPictureToAnimal(profilePicture: string): string  {
        const filename = profilePicture.split('/').pop() || ''; // Gets the last segment after '/'
        return animalMapping[filename] || 'Default'; 
    }

    return (
        <SafeAreaView  style={styles.container}>
            <Text style={styles.header}>GiveGetGo</Text>
            {/* <Text style={styles.userInfo}>Updated User Info: {JSON.stringify(userInfo, null, 2)}</Text> */ /*testing if the userInfo file got successfully updated*/} 
            <Avatar.Image source={getProfilePictureSource(userInfo.profilePicture)} style={styles.settings_avatar} /> 
            <Text style={styles.name}>{name}</Text>
            <Text style={styles.details}>{`Class of ${userInfo.classYear} • ${userInfo.major}`}</Text>
            <Text style={styles.details_bio}>{bio}</Text>
            <Divider style={styles.divider} />
            <Card style={styles.settings_card}>
                {/* <Card.Title title="Name" /> */}
                <Card.Content style={styles.settings_card_content}>
                    {isEditing_name ? (
                    <TextInput
                        style={styles.codeInput}
                        value={name}
                        onChangeText={setName}
                        returnKeyType="done"
                        right={<TextInput.Icon icon="content-save" onPress={handleSavePress_name} />}
                    />
                    ) : (
                    <List.Item
                        title="Name"
                        titleStyle={styles.settings_title}
                        descriptionStyle={styles.settings_description}
                        description={name}
                        right={props => 
                            <IconButton
                                icon="pencil"
                                size={20}
                                onPress={handleEditPress_name}
                            />}
                    />
                    )}
                </Card.Content>
            </Card>
            <Divider />
            <Card style={styles.settings_card}>
                {/* <Card.Title title="Email" /> */}
                <Card.Content style={styles.settings_card_content}>
                    {isEditing_email ? (
                    <TextInput
                        style={styles.codeInput}
                        value={email}
                        onChangeText={setEmail}
                        returnKeyType="done"
                        right={<TextInput.Icon icon="content-save" onPress={handleSavePress_email} />}
                    />
                    ) : (
                    <List.Item
                        title="Email"
                        titleStyle={styles.settings_title}
                        descriptionStyle={styles.settings_description}
                        description={email}
                        right={props => 
                            <IconButton
                                icon="pencil"
                                size={20}
                                onPress={handleEditPress_email}
                            />}
                    />
                    )}
                </Card.Content>
            </Card>
            <Divider />
            <Card style={styles.settings_card}>
                {/* <Card.Title title="Bio" /> */}
                <Card.Content style={styles.settings_card_content}>
                    {isEditing_bio ? (
                    <TextInput
                        style={styles.codeInput}
                        value={bio}
                        onChangeText={setBio}
                        returnKeyType="done"
                        right={<TextInput.Icon icon="content-save" onPress={handleSavePress_bio} />}
                    />
                    ) : (
                    <List.Item
                        title="Bio"
                        titleStyle={styles.settings_title}
                        descriptionStyle={styles.settings_description}
                        description={bio}
                        right={props => 
                            <IconButton
                                icon="pencil"
                                size={20}
                                onPress={handleEditPress_bio}
                            />}
                    />
                    )}
                </Card.Content>
            </Card>
            <Divider />
            <Card style={styles.settings_card}>
                {/* <Card.Title title="Profile Picture" /> */}
                <Card.Content style={styles.settings_card_content}>
                    {/* {isEditing_pic ? (
                    <TextInput
                        style={styles.codeInput}
                        value={profilePicture}
                        onChangeText={setProfilePicture}
                        returnKeyType="done"
                        right={<TextInput.Icon icon="content-save" onPress={handleSavePress_pic} />}
                    />
                    ) : (
                    <List.Item
                        title="Profile Picture"
                        titleStyle={styles.settings_title}
                        descriptionStyle={styles.settings_description}
                        description={profilePicture}
                        right={props => 
                            <IconButton
                                icon="pencil"
                                size={20}
                                onPress={handleEditPress_pic}
                            />}
                    />
                    )} */}
                    <List.Item
                        title="Profile Avatar"
                        titleStyle={styles.settings_title}
                        descriptionStyle={styles.settings_description}
                        description={mapPictureToAnimal(userInfo.profilePicture)}
                        right={props => 
                            <IconButton
                                icon="pencil"
                                size={20}
                                onPress={handleEditPress_pic}
                            />}
                    />
                </Card.Content>
            </Card>
            <Divider style={styles.divider} />
            <Button mode="contained" onPress={handleLogOut} style={styles.logOutButton}>
                Log Out
            </Button>
        </SafeAreaView >
    );
};

const styles = StyleSheet.create({
container: {
    flex: 1,
    marginTop: 50,
},
headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    justifyContent: 'space-between', // Distributes items evenly
    paddingLeft: 10, // Adds padding to the left of the avatar
    paddingRight: 10, // Adds padding to the right side
},
header: {
    fontSize: 22, // Increase the font size
    fontWeight: '600', // Make the font weight bold
    fontStyle: 'italic',
    textAlign: 'center', // Center the text
},
card: {
    margin: 8,
    overflow: 'scroll',
},
column: {
    justifyContent: 'space-between',
},
searchbar: {
    margin: 10,
},
resultsContainer: {
    marginTop: 20,
},
resultItem: {
    padding: 10,
    fontSize: 18,
},
filterButton: {
    position: 'absolute', 
    left: 280, 
    top: -20, 
    margin: 0, 
},
input: {
    marginBottom: 10,
},
settings_avatar: {
    alignSelf: 'center',
    marginTop: 20,
    backgroundColor: '#c7c7c7', // Placeholder color
},
name: {
    fontSize: 20,
    fontWeight: 'bold',
    textAlign: 'center',
    marginTop: 8,
},
details: {
    textAlign: 'center',
    marginBottom: 6,
},
details_bio: {
    textAlign: 'center',
    fontStyle: 'italic',
    marginBottom: 6,
},
divider: {
    marginVertical: 4,
},
logOutButton: {
    margin: 0,
},
settings_title: {
    fontSize: 16,
    fontWeight: 'bold',
},
settings_description: {
    fontSize: 14,
},
settings_card: {
    marginTop: 3,
    marginBottom: 3,
    maxHeight: 60,
},
settings_card_content: {
    marginTop: -10, 
    paddingTop: 0
},
notifications_card: {
    margin: 6,
    padding: 0,
    maxHeight: 150,
},
notifications_title: {
    fontSize: 16,
    fontWeight: 'bold',
    position: 'absolute', 
    top: 3, 
    left: 8, 
    padding: 8,
},
notifications_time: {
    fontSize: 12,
    color: '#757575', // grey color
    position: 'absolute', 
    top: 8, 
    right: 8, 
    padding: 8,
},
notifications_description: {
    fontSize: 15,
    top: 20, 
    right: 8, 
    padding: 9,
},
notifications_seeMore: {
    fontSize: 13,
    color: '#6200ee', // assuming a primary color
    fontWeight: 'bold',
    marginTop: -16,
},
deleteIcon: {
    width: 24, // Set the width of the icon
    height: 24, // Set the height of the icon
    alignSelf: 'center' // Center align the icon horizontally
  },
});

export default MainScreen;
