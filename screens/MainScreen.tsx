import React, { useState, useEffect } from 'react';
import { SafeAreaView, StyleSheet, View, FlatList, Dimensions, TouchableOpacity, TouchableWithoutFeedback, Keyboard, ScrollView } from 'react-native';
import { BottomNavigation, Text, Searchbar, Card, Title, Paragraph, Avatar, IconButton, Button, List, Divider, TextInput, Modal } from 'react-native-paper';
import { useIsFocused } from '@react-navigation/native';
import { NavigationContainer, CommonActions } from '@react-navigation/native';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import Icon from 'react-native-vector-icons/MaterialCommunityIcons';
import { NativeStackScreenProps, createNativeStackNavigator } from '@react-navigation/native-stack';
import { useRoute, RouteProp } from '@react-navigation/native';
import { useSelector, useDispatch } from 'react-redux';
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
import PostSubmittedScreen from './PostSubmittedScreen'; 
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';            
import * as Updates from 'expo-updates';
import { setAvatarUri } from '../store';

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
    PostRequestInfoScreen: { postId: string }; 
    PostRequestSucceedScreen: undefined;
    SettingsScreen: {newAvatarUri: string};
    AvatarPickerScreen: undefined;
    PostScreen: undefined;
    PostSubmittedScreen: undefined;
  };

type Post = {
    id: string;
    title?: string; // Optional, if your posts have titles
    description?: string; // Optional, if your posts have descriptions
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

type PostScreenProps = NativeStackScreenProps<RootStackParamList, 'PostScreen'>;

type NotificationsProps = NativeStackScreenProps<RootStackParamList, 'NotificationScreen'>;

type SettingsScreenProps = NativeStackScreenProps<RootStackParamList, 'SettingsScreen'>;

// dummy data for posts //should be loaded from the database
const posts_data: Post[] = [
    { id: '1', title: 'SCLA101 Homework Help', description: 'I need someone tutoring me for my essay assignment.'},
    { id: '2', title: 'ECE Career Coaching', description: 'I am seeking a career coach to help me set and achieve goals.'},
    { id: '3', title: 'Furniture Assembly', description: 'I need help assembling a new desk and bookshelf from IKEA.'},
    { id: '4', title: 'Meditation Guidance', description: 'Looking for a teacher to help me establish a regular meditation practice.'},
    { id: '5', title: 'Moving Assistance', description: 'I need help packing and moving my belongings to a new apartment across town.'},
    { id: '6', title: 'Need a Bike for Commuting', description: 'Looking for a used bicycle in good working condition to avoid walking 30 minutes to and from school every day.'},
    { id: '7', title: 'Taking Care of My Cat', description: 'I will be out of town on Febuary 28th and would need someone to look over my cat Lana.'},
    // ...more posts
  ];

// dummy data for notifications //should be loaded from the database
const notifications_data: Notification[] = [
    {
      id: '1',
      title: 'Request',
      time: '1 hour ago',
      description: 'New matching request from Rita Cheng for your "C Programming Guide".'
    },
    {
      id: '2',
      title: 'Match',
      time: '2 hours ago',
      description: 'Match succeeded with Jimmy Ho for "Python for Data Analysis". Click in to rate this match!'
    },
    {
      id: '3',
      title: 'Request',
      time: '3 hours ago',
      description: 'New matching request from Carol D. for your "Advanced Calculus".'
    },
    {
      id: '4',
      title: 'Match',
      time: '1 day ago',
      description: 'Match succeeded with Dave E. for "UX Design Fundamentals". Click in to rate this match!'
    },
    {
        id: '5',
        title: 'Match',
        time: '1 month ago',
        description: 'Match succeeded with Jason F. for "Crypto Trading 101". Click in to rate this match!'
      },
];

// dummy data for settings is written locally

// Custom hook to divide posts into two columns; used in HomeScreen
function useDividePostsIntoColumns(posts: Post[]): [Post[], Post[]] {
    const [columnOne, setColumnOne] = useState<Post[]>([]);
    const [columnTwo, setColumnTwo] = useState<Post[]>([]);

    useEffect(() => {
        const col1: Post[] = [];
        const col2: Post[] = [];

        posts.forEach((item, index) => {
            if (index % 2 === 0) col1.push(item);
            else col2.push(item);
        });

        setColumnOne(col1);
        setColumnTwo(col2);
    }, [posts]); // Dependency array includes posts to recalculate when posts change

    return [columnOne, columnTwo];
}

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

// Build Post's Stack
const PostStack = createNativeStackNavigator();
function PostStackScreen() {
  return (
    <PostStack.Navigator>
      <PostStack.Screen name="PostScreen" component={PostScreen} options={{ headerShown: false }}/>
      <PostStack.Screen name="PostSubmittedScreen" component={PostSubmittedScreen} options={{ headerShown: false }}/>
    </PostStack.Navigator>
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
                component={PostStackScreen}
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
    // const navigation = useNavigation<ScreenNavigationProp>();

    const [posts, setPosts] = React.useState<Post[]>(posts_data); 

    // useEffect(() => { //only reload after
 
    //     const fetchPostsData = async () => {
    //         try {
    //         const response = await fetch(`http://api.givegetgo.xyz/v1/post`, {
    //             method: 'GET',
    //             headers: {
    //             'Content-Type': 'application/json',
    //             },
    //         });

    //         if (response.ok) {
    //             const json = await response.json();
    //             setPosts(json.data); // Assuming the JSON response structure matches your state structure
    //         } else {
    //             throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    //         }
    //         } catch (error) {
    //         console.error('Error fetching posts:', error);
    //         // Optionally alert the user or handle the error visually in the UI
    //         }
    //     };

    // fetchPostsData();
    // }, []);

    const isFocused = useIsFocused();
    useEffect(() => {
        const fetchPostsData = async () => {
          try {
            const response = await fetch('http://api.givegetgo.xyz/v1/post', {
              method: 'GET',
              headers: {
                'Content-Type': 'application/json',
              },
            });
    
            if (response.ok) {
              const json = await response.json();
              setPosts(json.data); // Assuming the JSON response structure matches your state structure
            } else {
              throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
          } catch (error) {
            console.error('Error fetching posts:', error);
            // Optionally alert the user or handle the error visually in the UI
          }
        };
    
        if (isFocused) {
          fetchPostsData();
        }
    
      }, [isFocused]); 

    const handleProfilePress = () => {
        navigation.navigate('ProfileScreen'); 
    };

    const filteredPosts = posts.filter(post =>
        post.title?.toLowerCase().includes(searchQuery.toLowerCase())
      );
    
    const navigateToPostDetails = (postId: string) => {
        navigation.navigate('PostDetailsScreen', {postId: postId,});
    };

    // Substituting "flatlist" to allow cards in the same row not aligning each other
    const [columnOne, columnTwo] = useDividePostsIntoColumns(posts);

    const renderCard = (item: Post) => ( 
        <Card mode="elevated" style={[styles.home_card]} onPress={() => navigateToPostDetails(item.id)} >
            <Card.Content>
            <Title style={styles.homecard_title}>{item.title}</Title>
            <Paragraph style={styles.homecard_description}>{item.description}</Paragraph>
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
        <SafeAreaView  style={styles.home_container}>
            <View style={styles.home_headerContainer}>
                <TouchableOpacity onPress={handleProfilePress}>
                    <Avatar.Image size={40} source={getProfilePictureSource(selectedAvatarUri)} />
                </TouchableOpacity>
                <Text style={styles.home_header}>GiveGetGo</Text>
                <View style={styles.home_placeholder} />
            </View>
            <Searchbar
                placeholder="Search"
                onChangeText={onChangeSearch}
                value={searchQuery}
                style={styles.searchbar}
                // inputStyle={styles.searchbarInput}
                // placeholderTextColor="#FAFAFA" 
                // iconColor="#FAFAFA"   
            />
            <ScrollView contentContainerStyle={styles.twocol_container}>
                <View style={styles.column}>
                    {columnOne.map(item => renderCard(item))}
                </View>
                <View style={styles.column}>
                    {columnTwo.map(item => renderCard(item))}
                </View>
            </ScrollView>
        </SafeAreaView >
    );
};

const PostScreen: React.FC<PostScreenProps> = ({ navigation }: PostScreenProps) => {

  const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic }); 

  const handlePostPressApi = async () => {
    try {
        const response = await fetch('http://api.givegetgo.xyz/v1/post', {
            method: 'POST',
            credentials: "include",
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: title,
                description: description,
                category: 'gm'
              }),
        });

        const json = await response.json();
        console.log('Post submitted with Title:', title, 'and Description:', description);

        if (response.status === 201) {
            console.log('Post created successfully:', json);
            navigation.navigate('PostSubmittedScreen');
        } else {
            throw new Error(`Failed with status ${response.status}: ${json.msg}`);
        }
    } catch (error) {
        console.error('Error submitting post:', error);
        alert('Failed to submit post. Please try again.');
    }
};

  const handlePostPress = () => {
    handlePostPressApi()
    // navigation.navigate('PostSubmittedScreen')
    console.log('Post submitted with Title:', title, 'and Description:', description);
  };

  const [title, onChangeTitle] = React.useState<string>('');
  const [description, onChangeDescription] = React.useState<string>('');
    return (
        <TouchableWithoutFeedback onPress={Keyboard.dismiss} accessible={false}>
            <SafeAreaView style={styles.container}>
                <View style={styles.headerContainer}>
                    <View style={styles.backActionPlaceholder} />
                    <Text style={styles.header}>GiveGetGo</Text>
                    <View style={styles.backActionPlaceholder} />
                </View>
                <Card style={styles.post_card} mode="elevated">
                    <Card.Title title={( <Text style={styles.create_post_title_text}>Create a Post</Text> )} right={(props) => 
                        <Button {...props} style={styles.post_button} labelStyle={styles.post_button_text} onPress={handlePostPress}>Post</Button>} />
                    <Card.Content>
                        <TextInput
                            editable
                            style={styles.post_input_title}
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
                            style={styles.post_input_description}
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

    const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic }); 

    const [notifications, setNotifications] = React.useState<Notification[]>(notifications_data);      
    
    useEffect(() => {                                                                     //fill this in to get db info
        const fetchNotificationsData = async () => {
          try {
            const response = await fetch('URL_TO_YOUR_BACKEND/notifications_data_endpoint');
            const json = await response.json();
            setNotifications(json); // Adjust this depending on the structure of your JSON
          } catch (error) {
            // console.error(error); // uncomment this after finish frontend developing
          }
        };
        fetchNotificationsData();
      }, []);

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
            <Paragraph style={styles.notifications_dots}>...</Paragraph>
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
                <View style={styles.backActionPlaceholder} />
                <Text style={styles.header}>GiveGetGo</Text>
                <View style={styles.backActionPlaceholder} />
            </View>
            <FlatList
                data={notifications}
                style={styles.notifications_flatList}
                renderItem={renderItem}
                keyExtractor={(item) => item.id}
                // removeClippedSubviews={true} // If you want to improve performance on large lists, uncomment the line
            />
        </SafeAreaView >
    );
};

const SettingsScreen: React.FC<SettingsScreenProps> = ({ navigation }: SettingsScreenProps) => {  

    const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic }); 

    let currentUserInfo: UserInfo = { //should be loaded from the database
        name: 'Gilbert Hsu',
        email: 'xxx@gmail.com',
        bio: 'I am Gilbert the magician.',
        profilePicture: './profile_icon.jpg', // Placeholder for profile picture URI
        major: 'Computer Science',
        classYear: 2024,
      };

    const [userInfo, setUserInfo] = React.useState<UserInfo>(currentUserInfo);
    const dispatch = useDispatch(); 
    const isFocused = useIsFocused();
    useEffect(() => {
        const fetchUserInfo = async () => {
          try {
            const response = await fetch('http://api.givegetgo.xyz/v1/user/me', {
              method: 'GET',
              headers: {
                'Content-Type': 'application/json',
              },
            });
    
            if (response.ok) {
              const json = await response.json();
              console.log('Profile loaded successfully:', json);
              setUserInfo({
                name: json.data.username, 
                email: json.data.email,
                major: json.data.major,
                classYear: json.data.class,
                bio: json.data.profile_info, 
                profilePicture: json.data.profile_image
              });
              dispatch(setAvatarUri(json.data.profile_image));
              console.log("New Avatar from backend: ", json.data.profile_image)
              console.log("updated userInfo: ", userInfo)
            } else {
              // Handle errors
              console.error('Failed to fetch user info:', response.status);
            }
          } catch (error) {
            console.error('Network error when fetching user info:', error);
          }
        };
    
        if (isFocused) {
            fetchUserInfo();
          }
    }, [isFocused, dispatch]);

    const handleSaveProfile = async (updates: Partial<UserInfo>) => {
        try {
            const response = await fetch('http://api.givegetgo.xyz/v1/user/me', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: updates.name,
                    email: updates.email, // Assuming you can update email via PUT
                    class: updates.classYear,
                    major: updates.major,
                    profileInfo: updates.bio, // Assuming you want to update bio or similar
                    profile_image: updates.profilePicture
                }),
            });
    
            if (response.ok) {
                const json = await response.json();
                console.log('Profile updated successfully:', json);
                setUserInfo(updateUserInfo(userInfo, updates));
                // update redux
                dispatch(setAvatarUri(updates.profilePicture));
            } else {
                // Handle errors
                const errorJson = await response.json();
                alert(`Failed to update profile: ${errorJson.msg}`);
            }
        } catch (error) {
            console.error('Network error when updating user info:', error);
            alert('Failed to connect to the server. Please try again later.');
        }
    };

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

    // Effect to update local state when userInfo changes
    useEffect(() => {
        if (isFocused)
            setName(userInfo.name);
            setEmail(userInfo.email);
            setBio(userInfo.bio);
    }, [userInfo]);

    const handleEditPress_name = () => {
        setIsEditing_name(true);
    };

    const handleSavePress_name = () => {
        setUserInfo(prev => updateUserInfo(prev, { name: name })); //updated userInfo should also be saved to the database
        setIsEditing_name(false);
        handleSaveProfile({ name: name });
    };

    const handleEditPress_email = () => {
        setIsEditing_email(true);
    };

    const handleSavePress_email = () => {
        setUserInfo(prev => updateUserInfo(prev, { email: email }));
        setIsEditing_email(false);
        handleSaveProfile({ email: email });
    };

    const handleEditPress_bio = () => {
        setIsEditing_bio(true);
    };

    const handleSavePress_bio = () => {
        setUserInfo(prev => updateUserInfo(prev, { bio: bio }));
        setIsEditing_bio(false);
        handleSaveProfile({ bio: bio });
    };

    const handleEditPress_pic = () => {
        setIsEditing_pic(true);
        navigation.navigate('AvatarPickerScreen')  
    };

    const handleSavePress_profilePicture = (newAvatarUri:string) => {
        setUserInfo(prev => updateUserInfo(prev, { profilePicture: newAvatarUri }));
        handleSaveProfile({ profilePicture: newAvatarUri });
    };

    const handleLogOutApi = async () => {
        try {
            const response = await fetch('http://api.givegetgo.xyz/v1/user/logout', {
                method: 'GET',
                credentials: "include",
                headers: {
                    // 'Content-Type': 'application/json',
                },
            });
    
            const json = await response.json();
            console.log("Logout response:", json);
    
            if (response.status === 200) {
                console.log('User logged out successfully:', json);
                handleLogOut()
            } else if (response.status === 500) {
                console.error('Internal server error:', json.msg);
                alert(`Error: ${json.msg}`);
            } else {
                console.error('Unexpected error during logout:', json);
                alert(`Error: ${json.msg}`);
            }
        } catch (error) {
            console.error('Network error during logout:', error);
            alert('Failed to connect to the server. Please try again later.');
        }
    };

    // Restarting Expo
    async function restartApp() {
        try {
          await Updates.reloadAsync();
        } catch (error) {
          console.error('Failed to restart the app:', error);
        }
      }
      
    // Set up modal
    const [visible, setVisible] = useState(false);
    const showModal = () => setVisible(true);
    const hideModal = () => setVisible(false);
    
    const handleLogOut = () => {
        // add log out api here
        showModal(); // Show the modal
        setTimeout(() => {
            hideModal(); // Hide the modal
            restartApp(); // Then restart the app //does not work on web
        }, 4000); // Delay in milliseconds //used to be 2000
        // restartApp()
    };

    // Get newAvatarUri from AvatarPickerScreen
    type SettingsScreenRouteProp = RouteProp<RootStackParamList, 'SettingsScreen'>;
    const route = useRoute<SettingsScreenRouteProp>();
    const newAvatarUri = route.params?.newAvatarUri;
    console.log('New Avatar URI:', newAvatarUri);

    // Update profile picture when newAvatarUri changes
    useEffect(() => {
        if (isFocused)
            if (newAvatarUri) {
                console.log("newAvatarUri detected")
                // setUserInfo(prevUserInfo => ({
                //     ...prevUserInfo,
                //     profilePicture: newAvatarUri
                // }));
                handleSavePress_profilePicture(newAvatarUri)
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
            <View style={styles.headerContainer}>
                <View style={styles.backActionPlaceholder} />
                <Text style={styles.header}>GiveGetGo</Text>
                <View style={styles.backActionPlaceholder} />
            </View>
            <View style={styles.settingsContainer}>
                {/* <Text style={styles.userInfo}>Updated User Info: {JSON.stringify(userInfo, null, 2)}</Text> */ /*testing if the userInfo file got successfully updated*/} 
                <Avatar.Image source={getProfilePictureSource(userInfo.profilePicture)} size={80} style={styles.settings_avatar} /> 
                <Text style={styles.name}>{userInfo.name}</Text>
                <Text style={styles.details}>{`Class of ${userInfo.classYear} • ${userInfo.major}`}</Text>
                <Text style={styles.details_bio}>{userInfo.bio}</Text>
                <Divider style={styles.divider} />
                <Card style={styles.settings_card} mode="contained">
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
                            description={userInfo.name}
                            right={props => 
                                <IconButton
                                    icon="pencil"
                                    style={styles.settings_icon}
                                    size={20}
                                    onPress={handleEditPress_name}
                                />}
                        />
                        )}
                    </Card.Content>
                </Card>
                {/* <Divider /> */}
                <Card style={styles.settings_card} mode="contained">
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
                            description={userInfo.email}
                            right={props => 
                                <IconButton
                                    icon="pencil"
                                    style={styles.settings_icon}
                                    size={20}
                                    onPress={handleEditPress_email}
                                />}
                        />
                        )}
                    </Card.Content>
                </Card>
                {/* <Divider /> */}
                <Card style={styles.settings_card} mode="contained">
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
                            description={userInfo.bio}
                            right={props => 
                                <IconButton
                                    icon="pencil"
                                    style={styles.settings_icon}
                                    size={20}
                                    onPress={handleEditPress_bio}
                                />}
                        />
                        )}
                    </Card.Content>
                </Card>
                {/* <Divider /> */}
                <Card style={styles.settings_card} mode="contained">
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
                                    style={styles.settings_icon}
                                    size={20}
                                    onPress={handleEditPress_pic}
                                />}
                        />
                    </Card.Content>
                </Card>
                <Divider style={styles.divider} />
                <Button mode="contained" onPress={handleLogOutApi} style={styles.logOutButton}>
                    Log Out
                </Button>
                <Modal visible={visible} onDismiss={hideModal} contentContainerStyle={styles.modal}>
                    <Text style={styles.logOutText}>See you again soon!</Text>
                </Modal>
            </View>
        </SafeAreaView >
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
home_container: {
    flex: 1,
    marginTop: 50,
    justifyContent: 'center',
},
home_headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    justifyContent: 'space-between', // Distributes items evenly horizontally
    paddingLeft: 10, 
    paddingRight: 10, 
    position: 'absolute', // So that while setting card to the vertical middle, it still stays at the same place
    top: 0, 
    left: 6,
    right: 0,
    zIndex: 1, // Ensure the headerContainer is above the card
  },
home_header: {
    fontSize: 22, // Increase the font size
    fontWeight: '600', // Make the font weight bold
    fontFamily: 'Montserrat_700Bold_Italic',
    textAlign: 'center', // Center the text
    color: '#444444', // Dark gray color
},
home_placeholder: {
    width: 48, // This should match the width of the Appbar.BackAction for balance
    height: 52,
},
homecard_title: {
    fontSize: 16,
    fontWeight: '600',
    marginTop: 1,
    marginBottom: 1,
    lineHeight: 18,
    // color: '#FAFAFA',
},
homecard_description: {
    fontSize: 14,
    // fontStyle:'italic',
    marginBottom: 1,
    // color: '#FAFAFA',
},
home_card: { 
    borderRadius: 15, // Add rounded corners to the card
    // marginVertical: 6,
    margin: 8,
    justifyContent: 'center',
    padding: 3,
    // backgroundColor: 'black', 
},
twocol_container: {
    flexDirection: 'row',
    padding: 10,
},
column: {
    flex: 1,
},
searchbar: {
    margin: 10,
    marginTop: 65,   
    // backgroundColor: 'black',
},
// searchbarInput: {
//     color: '#FAFAFA',
// },
card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 12, // Add padding inside the card
    // marginTop: 170,
},
post_card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 15,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    padding: 12, // Add padding inside the card
    height: 400,
},
create_post_title_text: {
    fontSize: 18,
    textAlign: 'center',
    marginBottom: -5,
    marginRight: 100,
},
post_input_title: {
    padding: 0,
    fontSize: 16,
    marginBottom: 10, 
    marginTop: -5, 
},
post_input_description: {
    padding: 0,
    fontSize: 16,
    marginBottom: 5, 
    textAlignVertical: 'top',
    paddingBottom: 110,
},
titleButtonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 10,
},
post_button: {
    color: '#ffffff',
    paddingHorizontal: 12,
    paddingVertical: 8,
    marginRight: 5,
    fontSize: 18,
    elevation: 2,
},
post_button_text: {   
    fontSize: 15,          
},
resultsContainer: {
    marginTop: 20,
},
resultItem: {
    padding: 10,
    fontSize: 18,
},
settingsContainer: {
    marginTop: -15,
},
logOutButton: {
    position: 'absolute', 
    left: 10,
    right: 10, 
    bottom: -50,
    alignSelf: 'center', 
},
settings_avatar: {
    alignSelf: 'center',
    marginTop: 25,
    backgroundColor: '#c7c7c7', // Placeholder color
},
codeInput: {
    height: 40,
    marginVertical: 10,
},
name: {
    fontSize: 20,
    fontWeight: 'bold',
    textAlign: 'center',
    marginTop: 8,
    marginBottom: 4,
},
details: {
    textAlign: 'center',
    marginBottom: 4,
},
details_bio: {
    textAlign: 'center',
    fontStyle: 'italic',
    marginBottom: 6,
},
divider: {
    marginVertical: 8,
},
settings_title: {
    fontSize: 16,
    fontWeight: 'bold',
    marginTop: -3,
    marginBottom: 2,
},
settings_description: {
    fontSize: 14,
},
settings_card: {
    marginTop: 5,
    marginBottom: 5,
    maxHeight: 60,
    borderRadius: 15, // Add rounded corners to the card
    // marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 5, // Add padding inside the card
},
settings_card_content: {
    marginTop: -10, 
    paddingTop: 0
},
settings_icon: {
    marginTop: 2, 
    marginRight: -10, 
},
notifications_flatList: {
    marginTop: 55,
},
notifications_card: {
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 3, // Add padding inside the card
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
    top: 17, //20
    right: 8, 
    padding: 8,
},
notifications_dots: {
    top: 0, 
    flex: 1,
},
deleteIcon: {
    width: 24, // Set the width of the icon
    height: 24, // Set the height of the icon
    alignSelf: 'center', // Center align the icon horizontally
    marginTop: 1,
    marginBottom: 1,
  },
modal: {
    backgroundColor: 'white',
    padding: 20,
    margin: 20,
    borderRadius: 5,
    alignItems: 'center',
  },
logOutText: {
    fontSize: 16,
    marginBottom: 20,
    textAlign: 'center',
  },
});

export default MainScreen;
